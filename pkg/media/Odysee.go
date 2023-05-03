package media

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
	"golang.org/x/net/html"
)

type OdyseeAPI struct{}

func init() {
	OdyseeAPI := &OdyseeAPI{}
	MediaFactory.Register(OdyseeAPI)
}

type odyseeStruct struct {
	Context      string    `json:"@context"`
	Type         string    `json:"@type"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	UploadDate   time.Time `json:"uploadDate"`
	Duration     string    `json:"duration"`
	URL          string    `json:"url"`
	ContentURL   string    `json:"contentUrl"`
	EmbedURL     string    `json:"embedUrl"`
	Author       struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"author"`
	Thumbnail struct {
		Type string `json:"@type"`
		URL  string `json:"url"`
	} `json:"thumbnail"`
	Keywords        string `json:"keywords"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	PotentialAction struct {
		Type             string `json:"@type"`
		Target           string `json:"target"`
		StartOffsetInput string `json:"startOffset-input"`
	} `json:"potentialAction"`
}

func (sc *OdyseeAPI) getOdyseePage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func parse(text string) (data []string) {

	tkn := html.NewTokenizer(strings.NewReader(text))

	var vals []string

	var isLi bool

	for {

		tt := tkn.Next()

		switch {

		case tt == html.ErrorToken:
			return vals

		case tt == html.StartTagToken:

			t := tkn.Token()
			isLi = t.Data == "script"

		case tt == html.TextToken:

			t := tkn.Token()

			if isLi {
				vals = append(vals, t.Data)
			}

			isLi = false
		}
	}
}

func (sc *OdyseeAPI) GetMedia(url string, username string) []Media {
	data, _ := sc.getOdyseePage(url)
	parsedData := parse(string(data))[0]
	parsedData = strings.ReplaceAll(parsedData, "\n", "")
	parsedData = strings.ReplaceAll(parsedData, "\\", "")
	dataInBytes := []byte(parsedData)
	var tackData odyseeStruct
	err := json.Unmarshal(dataInBytes, &tackData)
	fmt.Println(err)

	media := []Media{}

	m := Media{
		ID:    ksuid.New().String(),
		Url:   url,
		User:  username,
		Type:  VIDEO_TYPE_ODYSEE,
		Title: tackData.Name,
		//Duration:    time.Duration(tackData.Duration / 1000),
		Thumbnail:   tackData.ThumbnailURL,
		ChannelName: tackData.Author.Name,
		AudioUrl:    tackData.ContentURL,
	}

	return append(media, m)
}

func (sc *OdyseeAPI) IsValidUrl(url string, ct *ContentType) bool {
	var re = regexp.MustCompile(`(?m)https:\/\/odysee.com\/(.*)`)
	match := re.Match([]byte(url))
	return match
}

func (sc *OdyseeAPI) GetType() string {
	return VIDEO_TYPE_ODYSEE
}
