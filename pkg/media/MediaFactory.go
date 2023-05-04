package media

import (
	"fmt"
	"net/http"
	"sort"
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
	GetType() string
	IsValidUrl(string, *ContentType) bool
	GetMedia(string, string) ([]Media, error)
	Refresh(media *Media) error
}

type Factory struct {
	Factories map[string]*MediaClient
}

var MediaFactory = Factory{
	Factories: map[string]*MediaClient{},
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

func (f *Factory) GetMedia(url string, username string) ([]Media, error) {
	factory := f.GetFactory(url)
	if factory == nil {
		return []Media{}, fmt.Errorf("unsupported URL")
	}
	return factory.GetMedia(url, username)
}

func NewVideo(url string, username string) ([]Media, error) {
	media, err := MediaFactory.GetMedia(url, username)
	return media, err
}

func RefreshAudioURL(media *Media) {
	factory := MediaFactory.GetFactory(media.Url)
	factory.Refresh(media)
}

func (f *Factory) GetTypes() []string {
	keys := make([]string, len(f.Factories))
	i := 0
	for k := range f.Factories {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
