package media

import "github.com/segmentio/ksuid"

type MP3Video struct{}

func init() {
	mp3Client := &MP3Video{}
	MediaFactory.Register(mp3Client)
}

func (client *MP3Video) GetMedia(url string, username string) []Media {
	return []Media{
		{
			ID:       ksuid.New().String(),
			Url:      url,
			User:     username,
			Type:     MediaType(client.GetType()),
			Title:    url,
			AudioUrl: url,
		},
	}
}

func (client *MP3Video) GetType() string {
	return VIDEO_TYPE_MP3
}

func (client *MP3Video) IsValidUrl(url string, ct *ContentType) bool {
	contentetType, err := ct.getConentType(url)
	if err != nil {
		return false
	}
	return contentetType == "audio/mpeg" || contentetType == "audio/aac" || contentetType == "audio/x-mpegurl"
}
