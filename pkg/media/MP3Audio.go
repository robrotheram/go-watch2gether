package media

import (
	"time"

	"github.com/segmentio/ksuid"
)

type MP3Video struct{}

func init() {
	mp3Client := &MP3Video{}
	MediaFactory.Register(mp3Client)
}

func (client *MP3Video) GetMedia(url string, username string) ([]Media, error) {
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

func (client *MP3Video) Refresh(media *Media) error {
	//Nothing to refresh
	return nil
}

func (client *MP3Video) GetType() MediaType {
	return VIDEO_TYPE_MP3
}

func (client *MP3Video) IsValidUrl(url string, ct *ContentType) bool {
	contentetType, err := ct.getConentType(url)
	if err != nil {
		return false
	}
	return contentetType == "audio/mpeg" || contentetType == "audio/aac" || contentetType == "audio/x-mpegurl"
}
