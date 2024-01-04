package media

import "github.com/segmentio/ksuid"

type MP4Video struct{}

func init() {
	mp4Client := &MP4Video{}
	MediaFactory.Register(mp4Client)
}

func (client *MP4Video) GetMedia(url string, username string) ([]Media, error) {
	return []Media{
		{
			ID:       ksuid.New().String(),
			Url:      url,
			User:     username,
			Type:     MediaType(client.GetType()),
			Title:    url,
			AudioUrl: url,
		},
	}, nil
}

func (client *MP4Video) Refresh(media *Media) error {
	//Nothing to refresh
	return nil
}

func (client *MP4Video) GetType() MediaType {
	return VIDEO_TYPE_MP4
}

func (client *MP4Video) IsValidUrl(url string, ct *ContentType) bool {
	contentetType, err := ct.getConentType(url)
	if err != nil {
		return false
	}
	return contentetType == "video/mp4" || contentetType == "application/dash+mpeg-url" || contentetType == "application/octet-stream" || contentetType == "application/vnd.apple.mpegurl"
}
