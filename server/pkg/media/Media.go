package media

import (
	"time"
)

type Media struct {
	ID        string        `json:"id"`
	User      string        `json:"user"`
	Url       string        `json:"url"`
	AudioUrl  string        `json:"audio_url"`
	Type      MediaType     `json:"type"`
	Title     string        `json:"title"`
	Channel   string        `json:channel`
	Duration  time.Duration `json:"duration"`
	Thumbnail string        `json:"thumbnail"`
	Order     int           `json:"order,omitempty"`
}

func (v *Media) GetType() MediaType {
	return v.Type
}

func (m *Media) GetAudioUrl() string {
	if m.AudioUrl == "" {
		return m.Url
	}
	return m.AudioUrl
}
