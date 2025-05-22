package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

type YTDLPMedia struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Formats []struct {
		FormatID   string  `json:"format_id"`
		FormatNote string  `json:"format_note,omitempty"`
		Ext        string  `json:"ext"`
		Protocol   string  `json:"protocol"`
		Acodec     string  `json:"acodec,omitempty"`
		Vcodec     string  `json:"vcodec"`
		URL        string  `json:"url"`
		Width      int     `json:"width,omitempty"`
		Height     int     `json:"height,omitempty"`
		Fps        float32 `json:"fps,omitempty"`
		Rows       int     `json:"rows,omitempty"`
		Columns    int     `json:"columns,omitempty"`
		Fragments  []struct {
			URL      string  `json:"url"`
			Duration float64 `json:"duration"`
		} `json:"fragments,omitempty"`
		AudioExt       string  `json:"audio_ext"`
		VideoExt       string  `json:"video_ext"`
		Vbr            float32 `json:"vbr"`
		Abr            float32 `json:"abr"`
		Tbr            any     `json:"tbr"`
		Resolution     string  `json:"resolution"`
		AspectRatio    float32 `json:"aspect_ratio"`
		FilesizeApprox any     `json:"filesize_approx,omitempty"`
		HTTPHeaders    struct {
			UserAgent      string `json:"User-Agent"`
			Accept         string `json:"Accept"`
			AcceptLanguage string `json:"Accept-Language"`
			SecFetchMode   string `json:"Sec-Fetch-Mode"`
		} `json:"http_headers"`
		Format             string  `json:"format"`
		FormatIndex        any     `json:"format_index,omitempty"`
		ManifestURL        string  `json:"manifest_url,omitempty"`
		Language           any     `json:"language,omitempty"`
		Preference         any     `json:"preference,omitempty"`
		Quality            float32 `json:"quality,omitempty"`
		HasDrm             bool    `json:"has_drm,omitempty"`
		SourcePreference   int     `json:"source_preference,omitempty"`
		Asr                float32 `json:"asr,omitempty"`
		Filesize           int     `json:"filesize,omitempty"`
		AudioChannels      int     `json:"audio_channels,omitempty"`
		LanguagePreference int     `json:"language_preference,omitempty"`
		DynamicRange       any     `json:"dynamic_range,omitempty"`
		Container          string  `json:"container,omitempty"`
		DownloaderOptions  struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options,omitempty"`
	} `json:"formats"`
	Thumbnails []struct {
		URL        string  `json:"url"`
		Preference int     `json:"preference"`
		ID         string  `json:"id"`
		Height     float32 `json:"height,omitempty"`
		Width      float32 `json:"width,omitempty"`
		Resolution string  `json:"resolution,omitempty"`
	} `json:"thumbnails"`
	Thumbnail         string   `json:"thumbnail"`
	Url               string   `json:"url"`
	Description       string   `json:"description"`
	ChannelID         string   `json:"channel_id"`
	ChannelURL        string   `json:"channel_url"`
	Duration          float32  `json:"duration"`
	ViewCount         float32  `json:"view_count"`
	AverageRating     any      `json:"average_rating"`
	AgeLimit          float32  `json:"age_limit"`
	WebpageURL        string   `json:"webpage_url"`
	Categories        []string `json:"categories"`
	Tags              []string `json:"tags"`
	PlayableInEmbed   bool     `json:"playable_in_embed"`
	LiveStatus        string   `json:"live_status"`
	ReleaseTimestamp  any      `json:"release_timestamp"`
	FormatSortFields  []string `json:"_format_sort_fields"`
	AutomaticCaptions struct {
	} `json:"automatic_captions"`
	Subtitles struct {
	} `json:"subtitles"`
	CommentCount int `json:"comment_count"`
	Chapters     any `json:"chapters"`
	Heatmap      []struct {
		StartTime float64 `json:"start_time"`
		EndTime   float64 `json:"end_time"`
		Value     float64 `json:"value"`
	} `json:"heatmap"`
	LikeCount            int     `json:"like_count"`
	Channel              string  `json:"channel"`
	ChannelFollowerCount int     `json:"channel_follower_count"`
	ChannelIsVerified    bool    `json:"channel_is_verified"`
	Uploader             string  `json:"uploader"`
	UploaderID           string  `json:"uploader_id"`
	UploaderURL          string  `json:"uploader_url"`
	UploadDate           string  `json:"upload_date"`
	Timestamp            float32 `json:"timestamp"`
	Availability         string  `json:"availability"`
	OriginalURL          string  `json:"original_url"`
	WebpageURLBasename   string  `json:"webpage_url_basename"`
	WebpageURLDomain     string  `json:"webpage_url_domain"`
	Extractor            string  `json:"extractor"`
	ExtractorKey         string  `json:"extractor_key"`
	Playlist             any     `json:"playlist"`
	PlaylistIndex        any     `json:"playlist_index"`
	DisplayID            string  `json:"display_id"`
	Fulltitle            string  `json:"fulltitle"`
	DurationString       string  `json:"duration_string"`
	ReleaseYear          any     `json:"release_year"`
	IsLive               bool    `json:"is_live"`
	WasLive              bool    `json:"was_live"`
	RequestedSubtitles   any     `json:"requested_subtitles"`
	HasDrm               any     `json:"_has_drm"`
	Epoch                int     `json:"epoch"`
	RequestedFormats     []struct {
		Asr                any     `json:"asr"`
		Filesize           int     `json:"filesize"`
		FormatID           string  `json:"format_id"`
		FormatNote         string  `json:"format_note"`
		SourcePreference   int     `json:"source_preference"`
		Fps                float32 `json:"fps"`
		AudioChannels      any     `json:"audio_channels"`
		Height             float32 `json:"height"`
		Quality            float64 `json:"quality"`
		HasDrm             bool    `json:"has_drm"`
		Tbr                float64 `json:"tbr"`
		FilesizeApprox     int     `json:"filesize_approx"`
		URL                string  `json:"url"`
		Width              float32 `json:"width"`
		Language           any     `json:"language"`
		LanguagePreference float32 `json:"language_preference"`
		Preference         any     `json:"preference"`
		Ext                string  `json:"ext"`
		Vcodec             string  `json:"vcodec"`
		Acodec             string  `json:"acodec"`
		DynamicRange       string  `json:"dynamic_range"`
		Container          string  `json:"container"`
		DownloaderOptions  struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options"`
		Protocol    string  `json:"protocol"`
		VideoExt    string  `json:"video_ext"`
		AudioExt    string  `json:"audio_ext"`
		Abr         float32 `json:"abr"`
		Vbr         float32 `json:"vbr"`
		Resolution  string  `json:"resolution"`
		AspectRatio float32 `json:"aspect_ratio"`
		HTTPHeaders struct {
			UserAgent      string `json:"User-Agent"`
			Accept         string `json:"Accept"`
			AcceptLanguage string `json:"Accept-Language"`
			SecFetchMode   string `json:"Sec-Fetch-Mode"`
		} `json:"http_headers"`
		Format string `json:"format"`
	} `json:"requested_formats"`
	Format         string  `json:"format"`
	FormatID       string  `json:"format_id"`
	Ext            string  `json:"ext"`
	Protocol       string  `json:"protocol"`
	Language       any     `json:"language"`
	FormatNote     string  `json:"format_note"`
	FilesizeApprox float32 `json:"filesize_approx"`
	Tbr            float32 `json:"tbr"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	Resolution     string  `json:"resolution"`
	Fps            float64 `json:"fps"`
	DynamicRange   string  `json:"dynamic_range"`
	Vcodec         string  `json:"vcodec"`
	Vbr            float64 `json:"vbr"`
	StretchedRatio any     `json:"stretched_ratio"`
	AspectRatio    float64 `json:"aspect_ratio"`
	Acodec         string  `json:"acodec"`
	Abr            float64 `json:"abr"`
	Asr            float32 `json:"asr"`
	AudioChannels  int     `json:"audio_channels"`
	Filename       string  `json:"_filename"`
	Filename0      string  `json:"filename"`
	Type           string  `json:"_type"`
	Version        struct {
		Version        string `json:"version"`
		CurrentGitHead any    `json:"current_git_head"`
		ReleaseGitHead any    `json:"release_git_head"`
		Repository     string `json:"repository"`
	} `json:"_version"`
}

type Client struct {
	Executable string
}

func processYTDLP(line string) (*Media, error) {

	var video YTDLPMedia
	err := json.Unmarshal([]byte(line), &video)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//Determin if the video is playable e.g if its has been privated or deleted
	if video.Title == "[Private video]" || video.Title == "[Deleted video]" {
		return nil, fmt.Errorf("video is unplayable")
	}

	if len(video.Thumbnail) == 0 && len(video.Thumbnails) > 0 {
		video.Thumbnail = video.Thumbnails[len(video.Thumbnails)-1].URL
	}

	media := &Media{
		ID:          ksuid.New().String(),
		Title:       video.Title,
		ChannelName: video.Channel,
		Url:         video.OriginalURL,
		Thumbnail:   video.Thumbnail,
		Progress: MediaDuration{
			Duration: time.Duration(video.Duration) * time.Second,
		},
		AudioUrl: video.Url,
	}

	return media, nil
}

func (client *Client) HydrateMedia(media *Media) error {
	cmd := exec.Command(client.Executable, "--flat-playlist", "--format", "worst", "--dump-json", media.Url)

	// Create a buffer to capture the output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return err
	}
	m, err := processYTDLP(out.String())
	if err != nil {
		return err
	}

	media.AudioUrl = m.AudioUrl

	return nil

}

func (client *Client) GetMedia(videoURL string, username string) ([]*Media, error) {
	var tracks = []*Media{}
	cmd := exec.Command(client.Executable, "--flat-playlist", "--dump-json", videoURL)

	// Create a buffer to capture the output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return tracks, err
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		media, err := processYTDLP(line)
		if err != nil {
			log.Warnf("unable to process audio %v", err)
			continue
		}
		media.AudioUrl = ""
		media.User = username
		tracks = append(tracks, media)
	}
	return tracks, nil
}
