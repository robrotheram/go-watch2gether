package media

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type MP4Video struct{}

func init() {
	mp4Client := &MP4Video{}
	MediaFactory.Register(mp4Client)
}

func getMediaInfo(url string) (*ffprobe.ProbeData, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFn()
	return ffprobe.ProbeURL(ctx, url)
}

func (client *MP4Video) GetMedia(url string, username string) ([]Media, error) {
	track := Media{
		ID:       ksuid.New().String(),
		Url:      url,
		User:     username,
		Type:     MediaType(client.GetType()),
		Title:    url,
		AudioUrl: url,
	}
	if data, err := getMediaInfo(url); err == nil {
		track.Title = data.Format.Filename
		track.Progress = MediaDuration{
			Duration: time.Duration(data.Format.DurationSeconds * float64(time.Second)),
		}
	}
	return []Media{track}, nil
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
	return contentetType == "video/x-matroska" || contentetType == "video/mp4" || contentetType == "application/dash+mpeg-url" || contentetType == "application/octet-stream" || contentetType == "application/vnd.apple.mpegurl"
}
