package media

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
)

var re = regexp.MustCompile(`(?m)"contentUrl": "(.*)"`)

type SoundCloudApi struct{}

// func init() {
// 	SoundCloudClient := &SoundCloudApi{}
// 	MediaFactory.Register(SoundCloudClient)
// }

type SoundcloudTrackInfo []struct {
	ArtworkURL        string      `json:"artwork_url"`
	Caption           interface{} `json:"caption"`
	Commentable       bool        `json:"commentable"`
	CommentCount      int         `json:"comment_count"`
	CreatedAt         time.Time   `json:"created_at"`
	Description       string      `json:"description"`
	Downloadable      bool        `json:"downloadable"`
	DownloadCount     int         `json:"download_count"`
	Duration          int         `json:"duration"`
	FullDuration      int         `json:"full_duration"`
	EmbeddableBy      string      `json:"embeddable_by"`
	Genre             string      `json:"genre"`
	HasDownloadsLeft  bool        `json:"has_downloads_left"`
	ID                int         `json:"id"`
	Kind              string      `json:"kind"`
	LabelName         interface{} `json:"label_name"`
	LastModified      time.Time   `json:"last_modified"`
	License           string      `json:"license"`
	LikesCount        int         `json:"likes_count"`
	Permalink         string      `json:"permalink"`
	PermalinkURL      string      `json:"permalink_url"`
	PlaybackCount     int         `json:"playback_count"`
	Public            bool        `json:"public"`
	PublisherMetadata struct {
		ID            int    `json:"id"`
		Urn           string `json:"urn"`
		ContainsMusic bool   `json:"contains_music"`
	} `json:"publisher_metadata"`
	PurchaseTitle interface{} `json:"purchase_title"`
	PurchaseURL   interface{} `json:"purchase_url"`
	ReleaseDate   interface{} `json:"release_date"`
	RepostsCount  int         `json:"reposts_count"`
	SecretToken   interface{} `json:"secret_token"`
	Sharing       string      `json:"sharing"`
	State         string      `json:"state"`
	Streamable    bool        `json:"streamable"`
	TagList       string      `json:"tag_list"`
	Title         string      `json:"title"`
	TrackFormat   string      `json:"track_format"`
	URI           string      `json:"uri"`
	Urn           string      `json:"urn"`
	UserID        int         `json:"user_id"`
	Visuals       interface{} `json:"visuals"`
	WaveformURL   string      `json:"waveform_url"`
	DisplayDate   time.Time   `json:"display_date"`
	Media         struct {
		Transcodings []struct {
			URL      string `json:"url"`
			Preset   string `json:"preset"`
			Duration int    `json:"duration"`
			Snipped  bool   `json:"snipped"`
			Format   struct {
				Protocol string `json:"protocol"`
				MimeType string `json:"mime_type"`
			} `json:"format"`
			Quality string `json:"quality"`
		} `json:"transcodings"`
	} `json:"media"`
	StationUrn         string `json:"station_urn"`
	StationPermalink   string `json:"station_permalink"`
	TrackAuthorization string `json:"track_authorization"`
	MonetizationModel  string `json:"monetization_model"`
	Policy             string `json:"policy"`
	User               struct {
		AvatarURL      string    `json:"avatar_url"`
		FirstName      string    `json:"first_name"`
		FollowersCount int       `json:"followers_count"`
		FullName       string    `json:"full_name"`
		ID             int       `json:"id"`
		Kind           string    `json:"kind"`
		LastModified   time.Time `json:"last_modified"`
		LastName       string    `json:"last_name"`
		Permalink      string    `json:"permalink"`
		PermalinkURL   string    `json:"permalink_url"`
		URI            string    `json:"uri"`
		Urn            string    `json:"urn"`
		Username       string    `json:"username"`
		Verified       bool      `json:"verified"`
		City           string    `json:"city"`
		CountryCode    string    `json:"country_code"`
		Badges         struct {
			Pro          bool `json:"pro"`
			ProUnlimited bool `json:"pro_unlimited"`
			Verified     bool `json:"verified"`
		} `json:"badges"`
		StationUrn       string `json:"station_urn"`
		StationPermalink string `json:"station_permalink"`
	} `json:"user"`
}

func (sc *SoundCloudApi) getSoundcloudData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (sc *SoundCloudApi) getTrackInfoUrl(html []byte) (string, error) {
	var re = regexp.MustCompile(`"transcodings":\[{"url":"(.+?)"`)
	trackID := re.FindAllStringSubmatch(string(html), -1)[0][1]
	return trackID, nil
}

func (sc *SoundCloudApi) getTrackID(html []byte) (string, error) {
	var re = regexp.MustCompile(`"soundcloud:tracks:(.+?)"`)
	trackID := re.FindAllStringSubmatch(string(html), -1)[0][1]
	return trackID, nil
}

func (sc *SoundCloudApi) getTrackData(trackID string, clientID string) SoundcloudTrackInfo {
	var track SoundcloudTrackInfo

	url := fmt.Sprintf("https://api-v2.soundcloud.com/tracks?ids=%s&client_id=%s", trackID, clientID)
	data, err := sc.getSoundcloudData(url)
	if err != nil {
		logrus.Warn("unable to get track data")
		return track
	}
	json.Unmarshal(data, &track)
	return track
}

func (sc *SoundCloudApi) getClientID(html []byte) (string, error) {
	var re = regexp.MustCompile(`<script crossorigin src="(.+?)"><\/script>`)
	url := re.FindAllStringSubmatch(string(html), -1)[4][1]

	data, err := sc.getSoundcloudData(url)
	if err != nil {
		return "", err
	}
	re = regexp.MustCompile(`client_id=(.+?)&`)
	clientID := re.FindAllStringSubmatch(string(data), -1)[0][1]

	return clientID, nil
}

func (sc *SoundCloudApi) GetMedia(url string, username string) []Media {
	media := []Media{}

	html, _ := sc.getSoundcloudData(url)
	trackInfo, err := sc.getTrackInfoUrl(html)
	if err != nil {
		return media
	}
	clientID := "iZIs9mchVcX5lhVRyQGGAYlNPVldzAoX"
	if err != nil {
		return media
	}
	trackID, err := sc.getTrackID(html)
	if err != nil {
		return media
	}
	tackData := sc.getTrackData(trackID, clientID)[0]

	type scMedia struct {
		Url string `json:"url"`
	}
	var scm scMedia

	data, err := sc.getSoundcloudData(fmt.Sprintf("%s?client_id=%s", trackInfo, clientID))
	if err != nil {
		return media
	}
	json.Unmarshal(data, &scm)

	m := Media{
		ID:          ksuid.New().String(),
		Url:         url,
		User:        username,
		Type:        VIDEO_TYPE_SC,
		Title:       tackData.Title,
		Duration:    time.Duration(tackData.Duration / 1000),
		Thumbnail:   tackData.ArtworkURL,
		ChannelName: tackData.User.Username,
		AudioUrl:    scm.Url,
	}

	return append(media, m)
}

func (sc *SoundCloudApi) IsValidUrl(url string, ct *ContentType) bool {
	re := regexp.MustCompile(`(?m)^((?:https?:)?\/\/)?((?:www|m)\.)?((?:soundcloud\.com))(\S+)?$`)
	match := re.Match([]byte(url))
	return match
}

func (sc *SoundCloudApi) GetType() string {
	return VIDEO_TYPE_SC
}
