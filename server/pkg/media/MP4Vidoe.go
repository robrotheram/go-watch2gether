package media

import "github.com/segmentio/ksuid"

type MP4Video struct{}

func init() {
	mp4Client := &MP4Video{}
	MediaFactory.Register(mp4Client)
}

func (client *MP4Video) GetMedia(url string, username string) []Media {
	return []Media{
		Media{
			ID:    ksuid.New().String(),
			Url:   url,
			User:  username,
			Type:  MediaType(client.GetType()),
			Title: url,
		},
	}
}

func (client *MP4Video) GetType() string {
	return VIDEO_TYPE_MP4
}

func (client *MP4Video) IsValidUrl(url string, ct *ContentType) bool {
	contentetType, err := ct.getConentType(url)
	if err != nil {
		return false
	}
	return contentetType == "video/mp4" || contentetType == "application/octet-stream" || contentetType == "application/vnd.apple.mpegurl"
}
