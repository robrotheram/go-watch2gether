package media

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/segmentio/ksuid"
	"golang.org/x/net/http/httpproxy"
)

var downloader *youtube.Client

func GetDownloader() *youtube.Client {

	if downloader != nil {
		return downloader
	}

	proxyFunc := httpproxy.FromEnvironment().ProxyFunc()
	httpTransport := &http.Transport{
		// Proxy: http.ProxyFromEnvironment() does not work. Why?
		Proxy: func(r *http.Request) (uri *url.URL, err error) {
			return proxyFunc(r.URL)
		},
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	downloader = &youtube.Client{}
	downloader.HTTPClient = &http.Client{Transport: httpTransport}

	return downloader
}

func GetVideoWithFormat(id string) (*youtube.Video, *youtube.Format, error) {
	dl := GetDownloader()
	video, err := dl.GetVideo(id)
	if err != nil {
		return nil, nil, err
	}
	formats := video.Formats

	if len(formats) == 0 {
		return nil, nil, errors.New("no formats found")
	}

	var format *youtube.Format
	// select the first format
	formats.Sort()
	format = &formats[len(formats)-1]
	return video, format, nil
}

func GetYoutubeURL(videoURL string) (string, error) {
	client := GetDownloader()
	video, format, err := GetVideoWithFormat(videoURL)
	if err != nil {
		return "", err
	}
	url, err := client.GetStreamURL(video, format)
	if err != nil {
		return "", err
	}
	return url, nil
}

func videosFromYoutubeURL(url string, username string) []Video {
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
				Type:     VIDEO_TYPE_YT,
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
			Type: VIDEO_TYPE_YT,
		}
		video.Update(ytVideo)
		return []Video{video}
	}
	return []Video{}
}