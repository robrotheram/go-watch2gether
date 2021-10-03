package media

import (
	"net/http"
	"regexp"
	"time"
)

type MediaType string

const (
	VIDEO_TYPE_YT      = "YOUTUBE"
	VIDEO_TYPE_MP4     = "MP4"
	VIDEO_TYPE_MP3     = "MP3"
	VIDEO_TYPE_PODCAST = "PODCAST"
)

func typeFromUrl(url string) MediaType {
	if isURLYT(url) {
		return VIDEO_TYPE_YT
	}
	contentType := ""
	resp, err := http.Head(url)
	if err != nil || resp.StatusCode != 200 {
		var client = &http.Client{
			Timeout: time.Second,
		}
		resp, err := client.Get(url)
		if err != nil {
			return ""
		}
		contentType = resp.Header.Get("Content-Type")
	} else {
		contentType = resp.Header.Get("Content-Type")
	}
	switch contentType {
	case "video/mp4":
		return VIDEO_TYPE_MP4
	case "audio/mpeg":
		return VIDEO_TYPE_MP3
	case "application/rss+xml; charset=UTF-8":
		return VIDEO_TYPE_PODCAST
	case "application/xml; charset=utf-8":
		return VIDEO_TYPE_PODCAST
	default:
		return ""
	}
}

func isURLYT(url string) bool {
	YT_REGEX := regexp.MustCompile(`(?m)^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)
	match := YT_REGEX.Match([]byte(url))
	return match
}
