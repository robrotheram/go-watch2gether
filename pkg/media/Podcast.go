package media

//https://github.com/nandosousafr/podfeed

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"
)

type Time struct {
	Value time.Time
}

func (t *Time) UnmarshalText(data []byte) (err error) {
	tm, err := time.Parse(time.RFC1123Z, string(data))
	if err != nil {
		tm, err = time.Parse(time.RFC1123, string(data))
		if err != nil {
			return
		}
	}
	*t = Time{tm}
	return
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.Value.String()), nil
}

type Podcast struct {
	Title       string   `xml:"channel>title"`
	Subtitle    string   `xml:"channel>subtitle"`
	Description string   `xml:"channel>description"`
	Link        string   `xml:"channel>link"`
	Language    string   `xml:"channel>language"`
	Author      string   `xml:"channel>author"`
	Image       Image    `xml:"channel>image"`
	Owner       Owner    `xml:"channel>owner"`
	Category    Category `xml:"channel>category"`
	Items       []Item   `xml:"channel>item"`
}

func (p Podcast) ReleasesByWeekday() (map[string]int, error) {
	res := map[string]int{}

	for _, episode := range p.Items {
		res[episode.PubDate.Value.Weekday().String()]++
	}

	return res, nil
}

type Item struct {
	Title       string    `xml:"title"`
	PubDate     Time      `xml:"pubDate"`
	Link        string    `xml:"link"`
	Duration    string    `xml:"duration"`
	Author      string    `xml:"author"`
	Summary     string    `xml:"summary"`
	Subtitle    string    `xml:"subtitle"`
	Description string    `xml:"description"`
	Enclosure   Enclosure `xml:"enclosure"`
	Image       Image     `xml:"image"`
}

type Image struct {
	Href  string `xml:"href,attr"`
	Url   string `xml:"url"`
	Title string `xml:"title"`
}

type Owner struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type Category struct {
	Text string `xml:"text,attr"`
}

type Enclosure struct {
	Type   string `xml:"type,attr"`
	Url    string `xml:"url,attr"`
	Length uint64 `xml:"length,attr"`
}
type PodcastRSS struct{}

func init() {
	PodcastClient := &PodcastRSS{}
	MediaFactory.Register(PodcastClient)
}

func (pod *PodcastRSS) parse(blob []byte) (pd Podcast, err error) {
	err = xml.Unmarshal(blob, &pd)
	if err != nil {
		return
	}

	return
}

func (pod *PodcastRSS) fetch(url string) (pd Podcast, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}

	defer res.Body.Close()

	buff, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return pod.parse(buff)
}

func (pod *PodcastRSS) GetMedia(url string, username string) ([]Media, error) {
	podcasts, err := pod.fetch(url)
	if err != nil {
		return []Media{}, err
	}

	item := podcasts.Items[0]

	return []Media{
		{
			Url:       item.Enclosure.Url,
			Type:      VIDEO_TYPE_PODCAST,
			Thumbnail: item.Image.Url,
			Title:     item.Title,
			User:      username,
			AudioUrl:  item.Enclosure.Url,
		},
	}, nil
}

func (pod *PodcastRSS) GetType() MediaType {
	return VIDEO_TYPE_PODCAST
}

func (pod *PodcastRSS) IsValidUrl(url string, ct *ContentType) bool {
	contentetType, err := ct.getConentType(url)
	if err != nil {
		return false
	}
	return contentetType == "application/rss+xml; charset=UTF-8" || contentetType == "application/xml; charset=utf-8"
}

func (client *PodcastRSS) Refresh(media *Media) error {
	//Nothing to refresh
	return nil
}
