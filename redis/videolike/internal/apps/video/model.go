package video

type likeVideoReq struct {
	UserId  int `json:"user_id" binding:"required,gte=1"`
	VideoId int `json:"video_id" binding:"required,gte=1"`
}
type likeVideoResp struct {
	VideoId        int `json:"video_id"`
	VideoLikeCount int `json:"video_like_count"`
}

type unlikeVideoReq struct {
	UserId  int `json:"user_id" binding:"required,gte=1"`
	VideoId int `json:"video_id" binding:"required,gte=1"`
}
type unlikeVideoResp struct {
	VideoId        int `json:"video_id"`
	VideoLikeCount int `json:"video_like_count"`
}
