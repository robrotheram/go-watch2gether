package media

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http/httpproxy"
)

type Youtube struct {
	downloader *youtube.Client
	hlsRegex   *regexp.Regexp
}

func init() {
	YoutubeClient := &Youtube{}
	YoutubeClient.Setup()
	MediaFactory.Register(YoutubeClient)
}

func (yt *Youtube) Setup() {

	if yt.downloader != nil {
		return
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

	yt.downloader = &youtube.Client{}
	yt.downloader.HTTPClient = &http.Client{Transport: httpTransport}
	yt.hlsRegex = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#,?&*//=]*)(.m3u8)\b([-a-zA-Z0-9@:%_\+.~#,?&//=]*))`)
}

func (yt *Youtube) GetVideoWithFormat(id string) (*youtube.Video, *youtube.Format, error) {
	video, err := yt.downloader.GetVideo(id)
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

func (yt *Youtube) getLiveURL(manifestURL string) (string, error) {

	res, err := http.Get(manifestURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if hlsURL := yt.hlsRegex.FindString(string(body)); hlsURL != "" {
		return hlsURL, nil
	} else {
		return "", errors.New("no valid URL found within HLS")
	}
}

func (yt *Youtube) GetAudioUrl(videoURL string) (string, error) {

	video, format, err := yt.GetVideoWithFormat(videoURL)
	if err != nil {
		return "", err
	}

	if video.HLSManifestURL != "" {
		logrus.Info("Getting HLS STREAM Url")
		url, err := yt.getLiveURL(video.HLSManifestURL)
		if err == nil {
			logrus.Info("USING HLS STREAM")
			return url, nil
		}
	}

	url, err := yt.downloader.GetStreamURL(video, format)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (yt *Youtube) GetMedia(url string, username string) []Media {
	ytPlaylist, err := yt.downloader.GetPlaylist(url)
	if err == nil {
		videos := []Media{}
		for _, ytVideo := range ytPlaylist.Videos {
			ytURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", ytVideo.ID)
			v := Media{
				ID:          ksuid.New().String(),
				Url:         ytURL,
				User:        username,
				Title:       ytVideo.Title,
				Type:        VIDEO_TYPE_YT,
				Duration:    ytVideo.Duration,
				ChannelName: ytVideo.Author,
			}
			if audio, err := yt.GetAudioUrl(ytURL); err == nil {
				v.AudioUrl = audio
			}
			videos = append(videos, v)
		}
		return videos
	}
	ytVideo, err := yt.downloader.GetVideo(url)
	if err == nil {
		video := Media{
			ID:          ksuid.New().String(),
			Url:         url,
			User:        username,
			Type:        VIDEO_TYPE_YT,
			Title:       ytVideo.Title,
			Duration:    ytVideo.Duration,
			Thumbnail:   ytVideo.Thumbnails[0].URL,
			ChannelName: ytVideo.Author,
		}
		if audio, err := yt.GetAudioUrl(url); err == nil {
			video.AudioUrl = audio
		}
		return []Media{video}
	}
	return []Media{}
}

func (yt *Youtube) GetType() string {
	return VIDEO_TYPE_YT
}

func (yt *Youtube) IsValidUrl(url string, ct *ContentType) bool {
	YT_REGEX := regexp.MustCompile(`(?m)^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)
	match := YT_REGEX.Match([]byte(url))
	return match
}
