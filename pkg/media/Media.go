package media

import (
	"time"
)

type Media struct {
	ID          string        `json:"id"`
	User        string        `json:"user"`
	Url         string        `json:"url"`
	AudioUrl    string        `json:"audio_url"`
	Type        MediaType     `json:"type"`
	Title       string        `json:"title"`
	ChannelName string        `json:"channel"`
	Progress    MediaDuration `json:"time"`
	Thumbnail   string        `json:"thumbnail"`
	Order       int           `json:"order,omitempty"`
}

type MediaDuration struct {
	Duration time.Duration `json:"duration"`
	Progress time.Duration `json:"progress"`
}

func (v *Media) GetType() MediaType {
	return v.Type
}

func (m *Media) GetAudioUrl() string {
	if m.AudioUrl == "" {
		m.Refresh()
	}
	return m.AudioUrl
}

func (m *Media) Refresh() error {
	return RefreshAudioURL(m)
}
