package media

import (
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/segmentio/ksuid"
)

type Video struct {
	ID        string        `json:"id"`
	User      string        `json:"user"`
	Url       string        `json:"url"`
	Type      MediaType     `json:"type"`
	Title     string        `json:"title"`
	Channel   string        `json:channel`
	Duration  time.Duration `json:"duration"`
	Thumbnail string        `json:"thumbnail"`
	Order     int           `json:"order,omitempty"`
}

func (v *Video) GetType() MediaType {
	return v.Type
}

func (v *Video) Update(m *youtube.Video) {
	v.Title = m.Title
	v.Duration = m.Duration
	v.Thumbnail = m.Thumbnails[0].URL
	v.Channel = m.Author
}

func NewVideo(url string, username string) ([]Video, error) {
	mediaType, err := TypeFromUrl(url)
	if err != nil {
		return []Video{}, err
	}

	switch mediaType {
	case VIDEO_TYPE_YT:
		return videosFromYoutubeURL(url, username), nil
	case VIDEO_TYPE_PODCAST:
		return videosFromPodcast(url, username), nil
	default:
		return []Video{
			{
				ID:    ksuid.New().String(),
				Url:   url,
				User:  username,
				Type:  mediaType,
				Title: url,
			},
		}, nil
	}
}
