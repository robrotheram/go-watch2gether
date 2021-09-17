package media

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/segmentio/ksuid"
)

const (
	VIDEO_TYPE_YT  = "YOUTUBE"
	VIDEO_TYPE_MP4 = "MP4"
)

var YT_REGEX = regexp.MustCompile(`(?m)^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)

type Video struct {
	ID        string        `json:"id"`
	User      string        `json:"user"`
	Url       string        `json:"url"`
	Type      string        `json:"type"`
	Title     string        `json:"title"`
	Channel   string        `json:channel`
	Duration  time.Duration `json:"duration"`
	Thumbnail string        `json:"thumbnail"`
	Order     int           `json:"order,omitempty"`
}

func (v *Video) GetType() string {
	if v.Type != "" {
		return v.Type
	}
	if YT_REGEX.Match([]byte(v.Url)) {
		v.Type = VIDEO_TYPE_YT
		return v.Type
	}
	if isURLMP4(v.Url) {
		v.Type = VIDEO_TYPE_MP4
		return v.Type
	}
	return ""
}

func isURLMP4(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	switch resp.Header.Get("Content-Type") {
	case "video/mp4":
		return true
	default:
		return false
	}
}

func (v *Video) Update(m *youtube.Video) {
	v.Title = m.Title
	v.Duration = m.Duration
	v.Thumbnail = m.Thumbnails[0].URL
	v.Channel = m.Author
}

func NewVideo(url string, username string) []Video {
	// Ignore non YT links
	if !YT_REGEX.Match([]byte(url)) {
		return []Video{
			{
				ID:   ksuid.New().String(),
				Url:  url,
				User: username,
			},
		}
	}
	client := GetDownloader()
	ytPlayist, err := client.GetPlaylist(url)
	if err == nil {
		vidoes := []Video{}
		for _, ytVideo := range ytPlayist.Videos {
			v := Video{
				ID:       ksuid.New().String(),
				Url:      fmt.Sprintf("https://www.youtube.com/watch?v=%s", ytVideo.ID),
				User:     username,
				Title:    ytVideo.Title,
				Duration: ytVideo.Duration,
				Channel:  ytVideo.Author,
			}
			vidoes = append(vidoes, v)
		}
		return vidoes
	}
	ytVideo, err := client.GetVideo(url)
	if err == nil {
		video := Video{
			ID:   ksuid.New().String(),
			Url:  url,
			User: username,
		}
		video.Update(ytVideo)
		return []Video{video}
	}
	return []Video{}
}
