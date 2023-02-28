package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func main() {
	var rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "redispass",
		DB:       0,
	})
	rdb.Ping(context.Background())

	var userId, videoId = 9, 10
	likeCount, err := likeVideo(rdb, userId, videoId)
	if err != nil {
		log.Fatalf("user(id=%d) likes video(id=%d) err", userId, videoId)
	}
	log.Printf("user(id=%d) likes video(id=%d), and video like count up to %d", userId, videoId, likeCount)

	likeCount, err = unlikeVideo(rdb, userId, videoId)
	if err != nil {
		log.Fatalf("user(id=%d) likes video(id=%d) err", userId, videoId)
	}
	log.Printf("user(id=%d) unlikes video(id=%d), and video like count up to %d", userId, videoId, likeCount)
}

func unlikeVideo(rdb *redis.Client, userId, videoId int) (videoLikeCount int, err error) {
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
	result, err := decrBy.Run(context.Background(), rdb, keys).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(result.(string))
}

func likeVideo(rdb *redis.Client, userId, videoId int) (videoLikeCount int, err error) {
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
	result, err := incrBy.Run(context.Background(), rdb, keys).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(result.(string))
}
