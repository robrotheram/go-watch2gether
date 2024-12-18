package media

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type ContentType struct {
	ContentType string
}

func (ct *ContentType) doRequest(_type string, url string) (*http.Response, error) {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	req, _ := http.NewRequest(_type, url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")
	return client.Do(req)
}

func (ct *ContentType) getConentType(url string) (string, error) {

	if ct.ContentType != "" {
		return ct.ContentType, nil
	}

	resp, err := ct.doRequest("HEAD", url)
	if err != nil || resp.StatusCode != 200 {
		resp, err := ct.doRequest("GET", url)
		if err != nil {
			return "", fmt.Errorf("unable to access url %v", err)
		}
		ct.ContentType = resp.Header.Get("Content-Type")
	} else {
		ct.ContentType = resp.Header.Get("Content-Type")
	}
	return ct.ContentType, nil
}

type MediaClient interface {
	GetType() MediaType
	IsValidUrl(string, *ContentType) bool
	GetMedia(url string, username string) ([]*Media, error)
	Refresh(media *Media) error
}

type Factory struct {
	Factories map[MediaType]*MediaClient
	Client    Client
}

var MediaFactory = Factory{
	Factories: map[MediaType]*MediaClient{},
	Client: Client{
		Executable: "yt-dlp",
	},
}

func (f *Factory) Register(client MediaClient) {
	f.Factories[client.GetType()] = &client
}

func (f *Factory) GetFactory(url string) MediaClient {
	ct := &ContentType{}
	for _, factory := range f.Factories {
		fact := *factory
		if fact.IsValidUrl(url, ct) {
			return *factory
		}
	}
	return nil
}

func (f *Factory) getMedia(url string, username string) ([]*Media, error) {
	if m, err := f.Client.GetMedia(url, username); err == nil {
		return m, err
	}
	factory := f.GetFactory(url)
	if factory == nil {
		return []*Media{}, fmt.Errorf("unsupported URL")
	}
	return factory.GetMedia(url, username)
}

func NewVideo(url string, username string) ([]*Media, error) {
	media, err := MediaFactory.getMedia(url, username)
	return media, err
}

func Refresh(media *Media) error {
	if len(media.AudioUrl) > 0 {
		log.Println("Skipping Refesh")
		return nil
	}
	if err := MediaFactory.Client.HydrateMedia(media); err == nil {
		return nil
	}
	factory := MediaFactory.GetFactory(media.Url)
	if factory == nil {
		return fmt.Errorf("no factory found")
	}
	return factory.Refresh(media)
}

func RefreshAll(tracks []*Media) error {
	for _, track := range tracks {
		Refresh(track)
	}
	return nil
}

func (f *Factory) GetTypes() []MediaType {
	keys := make([]MediaType, len(f.Factories))
	i := 0
	for k := range f.Factories {
		keys[i] = k
		i++
	}
	return keys
}
