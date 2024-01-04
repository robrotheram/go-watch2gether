package media

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/segmentio/ksuid"
)

var PeertubeServerRegex = regexp.MustCompile(`(?m)^(((https:)\/\/).*)\/w\/`)
var PeertubeVideoRegex = regexp.MustCompile(`(?m)/w\/(.*)`)

type PeertubeFile struct {
	ID         int `json:"id"`
	Resolution struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	} `json:"resolution"`
	MagnetURI          string `json:"magnetUri"`
	Size               int    `json:"size"`
	Fps                int    `json:"fps"`
	TorrentURL         string `json:"torrentUrl"`
	TorrentDownloadURL string `json:"torrentDownloadUrl"`
	FileURL            string `json:"fileUrl"`
	FileDownloadURL    string `json:"fileDownloadUrl"`
	MetadataURL        string `json:"metadataUrl"`
}
type PeertubeInfo struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	ShortUUID string `json:"shortUUID"`
	URL       string `json:"url"`
	Name      string `json:"name"`
	Category  struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	} `json:"category"`
	Licence struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	} `json:"licence"`
	Language struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"language"`
	Privacy struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	} `json:"privacy"`
	Nsfw                  bool        `json:"nsfw"`
	Description           string      `json:"description"`
	IsLocal               bool        `json:"isLocal"`
	Duration              int         `json:"duration"`
	Views                 int         `json:"views"`
	Viewers               int         `json:"viewers"`
	Likes                 int         `json:"likes"`
	Dislikes              int         `json:"dislikes"`
	ThumbnailPath         string      `json:"thumbnailPath"`
	PreviewPath           string      `json:"previewPath"`
	EmbedPath             string      `json:"embedPath"`
	CreatedAt             time.Time   `json:"createdAt"`
	UpdatedAt             time.Time   `json:"updatedAt"`
	PublishedAt           time.Time   `json:"publishedAt"`
	OriginallyPublishedAt interface{} `json:"originallyPublishedAt"`
	IsLive                bool        `json:"isLive"`
	Account               struct {
		URL     string `json:"url"`
		Name    string `json:"name"`
		Host    string `json:"host"`
		Avatars []struct {
			Width     int       `json:"width"`
			Path      string    `json:"path"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"avatars"`
		Avatar struct {
			Width     int       `json:"width"`
			Path      string    `json:"path"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"avatar"`
		ID                    int           `json:"id"`
		HostRedundancyAllowed bool          `json:"hostRedundancyAllowed"`
		FollowingCount        int           `json:"followingCount"`
		FollowersCount        int           `json:"followersCount"`
		CreatedAt             time.Time     `json:"createdAt"`
		Banners               []interface{} `json:"banners"`
		DisplayName           string        `json:"displayName"`
		Description           string        `json:"description"`
		UpdatedAt             time.Time     `json:"updatedAt"`
		UserID                int           `json:"userId"`
	} `json:"account"`
	Channel struct {
		URL     string `json:"url"`
		Name    string `json:"name"`
		Host    string `json:"host"`
		Avatars []struct {
			Width     int       `json:"width"`
			Path      string    `json:"path"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"avatars"`
		Avatar struct {
			Width     int       `json:"width"`
			Path      string    `json:"path"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"avatar"`
		ID                    int           `json:"id"`
		HostRedundancyAllowed bool          `json:"hostRedundancyAllowed"`
		FollowingCount        int           `json:"followingCount"`
		FollowersCount        int           `json:"followersCount"`
		CreatedAt             time.Time     `json:"createdAt"`
		Banners               []interface{} `json:"banners"`
		DisplayName           string        `json:"displayName"`
		Description           interface{}   `json:"description"`
		Support               interface{}   `json:"support"`
		IsLocal               bool          `json:"isLocal"`
		UpdatedAt             time.Time     `json:"updatedAt"`
		OwnerAccount          struct {
			URL     string `json:"url"`
			Name    string `json:"name"`
			Host    string `json:"host"`
			Avatars []struct {
				Width     int       `json:"width"`
				Path      string    `json:"path"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
			} `json:"avatars"`
			Avatar struct {
				Width     int       `json:"width"`
				Path      string    `json:"path"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
			} `json:"avatar"`
			ID                    int           `json:"id"`
			HostRedundancyAllowed bool          `json:"hostRedundancyAllowed"`
			FollowingCount        int           `json:"followingCount"`
			FollowersCount        int           `json:"followersCount"`
			CreatedAt             time.Time     `json:"createdAt"`
			Banners               []interface{} `json:"banners"`
			DisplayName           string        `json:"displayName"`
			Description           string        `json:"description"`
			UpdatedAt             time.Time     `json:"updatedAt"`
			UserID                int           `json:"userId"`
		} `json:"ownerAccount"`
	} `json:"channel"`
	Blacklisted        bool        `json:"blacklisted"`
	BlacklistedReason  interface{} `json:"blacklistedReason"`
	StreamingPlaylists []struct {
		ID                int            `json:"id"`
		Type              int            `json:"type"`
		PlaylistURL       string         `json:"playlistUrl"`
		SegmentsSha256URL string         `json:"segmentsSha256Url"`
		Redundancies      []interface{}  `json:"redundancies"`
		Files             []PeertubeFile `json:"files"`
	} `json:"streamingPlaylists"`
	Files           []PeertubeFile `json:"files"`
	Support         interface{}    `json:"support"`
	DescriptionPath string         `json:"descriptionPath"`
	Tags            []interface{}  `json:"tags"`
	CommentsEnabled bool           `json:"commentsEnabled"`
	DownloadEnabled bool           `json:"downloadEnabled"`
	WaitTranscoding bool           `json:"waitTranscoding"`
	State           struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	} `json:"state"`
	TrackerUrls []string `json:"trackerUrls"`
}

type Peertube struct {
	client *http.Client
}

func init() {
	PeertubeClient := &Peertube{client: &http.Client{Timeout: 10 * time.Second}}
	MediaFactory.Register(PeertubeClient)
}

func (pt *Peertube) getPeerTubeData(apiURL string) (PeertubeInfo, error) {
	var info PeertubeInfo
	r, err := pt.client.Get(apiURL)
	if err != nil {
		return info, err
	}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&info)
	return info, err
}

func (pt *Peertube) GetMedia(url string, username string) ([]Media, error) {
	media := []Media{}
	host := PeertubeServerRegex.FindStringSubmatch(url)
	if len(host) < 2 {
		return media, fmt.Errorf("unable to parse peertube url")
	}
	video := PeertubeVideoRegex.FindStringSubmatch(url)
	if len(video) < 2 {
		return media, fmt.Errorf("unable to parse peertube url")
	}
	apiURL := fmt.Sprintf("%s/api/v1/videos/%s", host[1], video[1])
	info, _ := pt.getPeerTubeData(apiURL)

	audioUrl := ""

	if len(info.Files) > 0 {
		for _, file := range info.Files {
			if file.Resolution.Label == "Audio" {
				audioUrl = file.FileURL
			}
		}
	} else if len(info.StreamingPlaylists) > 0 {
		//audioUrl = info.StreamingPlaylists[0].PlaylistURL
		playlist := info.StreamingPlaylists[0].Files[0]
		for _, play := range info.StreamingPlaylists {
			for _, file := range play.Files {
				if file.Size < playlist.Size {
					playlist = file
				}
			}
		}
		audioUrl = playlist.FileDownloadURL
	}

	m := Media{
		ID:    ksuid.New().String(),
		Url:   info.URL,
		User:  username,
		Type:  VIDEO_TYPE_PEERTUBE,
		Title: info.Name,
		Progress: MediaDuration{
			Duration: 0,
		},
		Thumbnail: fmt.Sprintf("%s%s", host[1], info.ThumbnailPath),
		AudioUrl:  audioUrl,
	}

	return append(media, m), nil
}

func (pt *Peertube) IsValidUrl(url string, ct *ContentType) bool {
	re := regexp.MustCompile(`(https:\/\/)?(?:.*)\/w\/(\S+)?$`)
	match := re.Match([]byte(url))
	return match
}

func (pt *Peertube) GetType() MediaType {
	return VIDEO_TYPE_PEERTUBE
}

func (client *Peertube) Refresh(media *Media) error {
	//Nothing to refresh
	return nil
}
