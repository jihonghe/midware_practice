package video

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	rediscli "pkg/redis"
)

func like(ctx *gin.Context) {
	var err error
	var req likeVideoReq
	var resp likeVideoResp

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("parse req body err: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "parse request body err"})
		return
	}

	resp.VideoId = req.VideoId
	resp.VideoLikeCount, err = likeVideo(req.UserId, req.VideoId)
	if err != nil {
		log.Printf("like video in redis err: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": fmt.Sprintf("like video(id=%d) err", req.VideoId)})
		return
	}

	err = myRepo.likeVideo(req.VideoId)
	if err != nil {
		log.Printf("like video in repo err: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": fmt.Sprintf("like video(id=%d) err", req.VideoId)})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func unlike(ctx *gin.Context) {
	var err error
	var req unlikeVideoReq
	var resp unlikeVideoResp

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "parse request body err"})
		return
	}

	resp.VideoId = req.VideoId
	resp.VideoLikeCount, err = unlikeVideo(req.UserId, req.VideoId)
	if err != nil {
		log.Printf("unlike video err: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "unlike video err"})
		return
	}

	err = myRepo.unlikeVideo(req.VideoId)
	if err != nil {
		log.Printf("unlike video in repo err: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": fmt.Sprintf("unlike video(id=%d) err", req.VideoId)})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func unlikeVideo(userId, videoId int) (videoLikeCount int, err error) {
	var decrBy = redis.NewScript(`
--取消点赞操作，删除并递减，如果删除失败则不递减，【原子性、幂等性】
if redis.call('DEL',KEYS[1]) == 1
then
    redis.call('DECR',KEYS[2])
end
return redis.call('GET',KEYS[2])
`)
	var keys = []string{
		fmt.Sprintf("u.%d.like.%d", userId, videoId),
		fmt.Sprintf("v.%d.like_count", videoId),
	}
	result, err := decrBy.Run(context.Background(), rediscli.Client, keys).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(result.(string))
}

func likeVideo(userId, videoId int) (videoLikeCount int, err error) {
	var incrBy = redis.NewScript(`
--点赞操作，写入并自增，如果写入失败则不自增，【原子性、幂等性】
if redis.call('SETNX',KEYS[1],1) == 1
then
    redis.call('EXPIRE',KEYS[1],864000)
    redis.call('INCR',KEYS[2])
end
return redis.call('GET',KEYS[2])
`)
	var keys = []string{
		fmt.Sprintf("u.%d.like.%d", userId, videoId),
		fmt.Sprintf("v.%d.like_count", videoId),
	}
	result, err := incrBy.Run(context.Background(), rediscli.Client, keys).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(result.(string))
}
