package media

import (
	"fmt"
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

func doRequest(_type string, url string) (*http.Response, error) {
	var client = &http.Client{
		Timeout: time.Second,
	}
	req, _ := http.NewRequest(_type, url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")
	return client.Do(req)
}

func typeFromUrl(url string) (MediaType, error) {
	if isURLYT(url) {
		return VIDEO_TYPE_YT, nil
	}
	contentType := ""
	resp, err := doRequest("HEAD", url)
	if err != nil || resp.StatusCode != 200 {
		resp, err := doRequest("GET", url)
		if err != nil {
			return "", fmt.Errorf("unable to access url")
		}
		contentType = resp.Header.Get("Content-Type")
	} else {
		contentType = resp.Header.Get("Content-Type")
	}
	switch contentType {
	case "video/mp4":
		return VIDEO_TYPE_MP4, nil
	case "audio/mpeg":
		return VIDEO_TYPE_MP3, nil
	case "audio/aac":
		return VIDEO_TYPE_MP3, nil
	case "application/rss+xml; charset=UTF-8":
		return VIDEO_TYPE_PODCAST, nil
	case "application/xml; charset=utf-8":
		return VIDEO_TYPE_PODCAST, nil
	default:
		return "", fmt.Errorf("unsupported media type %s", contentType)
	}
}

func isURLYT(url string) bool {
	YT_REGEX := regexp.MustCompile(`(?m)^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)
	match := YT_REGEX.Match([]byte(url))
	return match
}
