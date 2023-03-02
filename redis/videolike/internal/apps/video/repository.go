package video

var (
	myRepo repo
)

type repo struct{}

func (r *repo) likeVideo(vid int) (err error) {
	return nil
}

func (r *repo) unlikeVideo(vid int) (err error) {
	return nil
}
