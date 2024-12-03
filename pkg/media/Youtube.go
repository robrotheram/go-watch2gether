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
	"w2g/pkg/media/clients/youtube"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
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

	yt.downloader = &youtube.Client{
		HTTPClient: &http.Client{Transport: httpTransport},
	}
	yt.downloader.SetDefaultClient(&youtube.WebClient)
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
	for _, f := range formats {
		if f.AudioChannels > 0 {
			return video, &f, nil
		}
	}
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
		log.Info("Getting HLS STREAM Url")
		url, err := yt.getLiveURL(video.HLSManifestURL)
		if err == nil {
			log.Info("USING HLS STREAM")
			return url, nil
		}
	}

	url, err := yt.downloader.GetStreamURL(video, format)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (yt *Youtube) GetMedia(url string, username string) ([]Media, error) {
	ytPlaylist, err := yt.downloader.GetPlaylist(url)
	if err == nil {
		videos := []Media{}
		for _, ytVideo := range ytPlaylist.Videos {
			ytURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", ytVideo.ID)

			v := Media{
				ID:    ksuid.New().String(),
				Url:   ytURL,
				User:  username,
				Title: ytVideo.Title,
				Type:  VIDEO_TYPE_YT,
				Progress: MediaDuration{
					Duration: ytVideo.Duration,
				},
				ChannelName: ytVideo.Author,
				Thumbnail:   ytVideo.Thumbnails[0].URL,
			}
			videos = append(videos, v)
		}
		return videos, nil
	}
	ytVideo, err := yt.downloader.GetVideo(url)

	if err == nil {
		video := Media{
			ID:          ksuid.New().String(),
			Url:         url,
			User:        username,
			Type:        isLive(ytVideo),
			Title:       ytVideo.Title,
			Progress:    getDuration(ytVideo),
			Thumbnail:   ytVideo.Thumbnails[0].URL,
			ChannelName: ytVideo.Author,
		}
		return []Media{video}, nil
	}
	return []Media{}, fmt.Errorf("unable find valid audio for this url, %v", err)
}

func (yt *Youtube) GetType() MediaType {
	return VIDEO_TYPE_YT
}

func (yt *Youtube) IsValidUrl(url string, ct *ContentType) bool {
	YT_REGEX := regexp.MustCompile(`(?m)^((?:https?:)?\/\/)?((?:www|m|music)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)
	match := YT_REGEX.Match([]byte(url))
	return match
}

func (yt *Youtube) Refresh(media *Media) error {
	audio, err := yt.GetAudioUrl(media.Url)
	media.AudioUrl = audio
	return err
}

func isLive(yt *youtube.Video) MediaType {
	if len(yt.HLSManifestURL) > 1 {
		return VIDEO_TYPE_YT_LIVE
	}
	return VIDEO_TYPE_YT
}

func getDuration(ytVideo *youtube.Video) MediaDuration {
	if len(ytVideo.HLSManifestURL) > 1 {
		return MediaDuration{
			Duration: -1,
		}
	}
	return MediaDuration{
		Duration: ytVideo.Duration,
	}
}
