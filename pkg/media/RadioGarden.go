package media

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"

	"github.com/segmentio/ksuid"
)

type RadioGardenInfo struct {
	APIVersion int    `json:"apiVersion"`
	Version    string `json:"version"`
	Data       struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		ID      string `json:"id"`
		URL     string `json:"url"`
		Website string `json:"website"`
		Secure  bool   `json:"secure"`
		Place   struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"place"`
		Country struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"country"`
	} `json:"data"`
}

type RadioGarden struct {
	client *http.Client
}

func init() {
	RadioGardenClient := &RadioGarden{client: &http.Client{Timeout: 10 * time.Second}}
	MediaFactory.Register(RadioGardenClient)
}

func (radio *RadioGarden) getRadioInfo(id string) (RadioGardenInfo, error) {
	var info RadioGardenInfo
	r, err := radio.client.Get(fmt.Sprintf("https://radio.garden/api/ara/content/channel/%s", id))
	if err != nil {
		return info, err
	}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&info)
	return info, err
}

func (radio *RadioGarden) getAudioURL(id string) string {
	return fmt.Sprintf("http://radio.garden/api/ara/content/listen/%s/channel.mp3", id)
}

func (radio *RadioGarden) getIDFromURL(targetUrl string) string {
	myUrl, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatal(err)
	}
	return (path.Base(myUrl.Path))
}
func (radio *RadioGarden) GetMedia(url string, username string) ([]Media, error) {
	media := []Media{}

	radioID := radio.getIDFromURL(url)

	info, _ := radio.getRadioInfo(radioID)
	audioUrl := radio.getAudioURL(radioID)

	m := Media{
		ID:        ksuid.New().String(),
		Url:       url,
		User:      username,
		Type:      VIDEO_TYPE_RG,
		Title:     info.Data.Title,
		Duration:  time.Duration(0),
		Thumbnail: "https://play-lh.googleusercontent.com/07lewhVI4GklVBi_ehhOXxmB_bPaWWTiyqHAlQP6VsYD7h9R4d8hskNAy4SCOx0leNx-=s180",
		AudioUrl:  audioUrl,
	}

	return append(media, m), nil
}

func (radio *RadioGarden) IsValidUrl(url string, ct *ContentType) bool {
	re := regexp.MustCompile(`(?m)^((?:https?:)?\/\/)?((?:www|m)\.)?((?:radio\.garden))(\S+)?$`)
	match := re.Match([]byte(url))
	return match
}

func (client *RadioGarden) GetType() string {
	return VIDEO_TYPE_RG
}

func (client *RadioGarden) Refresh(media *Media) error {
	//Nothing to refresh
	return nil
}
