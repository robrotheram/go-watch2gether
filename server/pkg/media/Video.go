package media

import (
	"net/http"
	"regexp"
)

const (
	VIDEO_TYPE_YT  = "YOUTUBE"
	VIDEO_TYPE_MP4 = "MP4"
)

var YT_REGEX = regexp.MustCompile(`(?m)^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)

type Video struct {
	ID    string `json:"id"`
	User  string `json:"user"`
	Url   string `json:"url"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Order int    `json:"order,omitempty"`
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
