package media

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
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
	return io.ReadAll(resp.Body)
}

func parseDuration(durationStr string) (time.Duration, error) {
	re := regexp.MustCompile(`PT(\d+)M(\d+)S`)
	matches := re.FindStringSubmatch(durationStr)
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	minutes, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, err
	}

	return time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
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

func (sc *OdyseeAPI) GetMedia(url string, username string) ([]Media, error) {
	data, _ := sc.getOdyseePage(url)
	scripts := parse(string(data))

	var tackData odyseeStruct
	for _, script := range scripts {
		script = strings.ReplaceAll(script, "\n", "")
		script = strings.ReplaceAll(script, "\\", "")
		dataInBytes := []byte(script)
		if json.Unmarshal(dataInBytes, &tackData) == nil {
			break
		}
	}
	media := []Media{}
	if tackData.ContentURL == "" {
		return media, fmt.Errorf("unable to parse site")
	}

	m := Media{
		ID:          ksuid.New().String(),
		Url:         url,
		User:        username,
		Type:        VIDEO_TYPE_ODYSEE,
		Title:       tackData.Name,
		Thumbnail:   tackData.ThumbnailURL,
		ChannelName: tackData.Author.Name,
		AudioUrl:    tackData.ContentURL,
	}
	if duration, err := parseDuration(tackData.Duration); err == nil {
		m.Progress = MediaDuration{
			Duration: duration,
		}
	}
	return append(media, m), nil
}

func (sc *OdyseeAPI) IsValidUrl(url string, ct *ContentType) bool {
	var re = regexp.MustCompile(`(?m)https:\/\/odysee.com\/(.*)`)
	match := re.Match([]byte(url))
	return match
}

func (sc *OdyseeAPI) GetType() MediaType {
	return VIDEO_TYPE_ODYSEE
}

func (client *OdyseeAPI) Refresh(media *Media) error {
	//Nothing to refresh
	return nil
}
