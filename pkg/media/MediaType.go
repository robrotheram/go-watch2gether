package media

type MediaType string

const (
	VIDEO_TYPE_YT       = MediaType("YOUTUBE")
	VIDEO_TYPE_YT_LIVE  = MediaType("YOUTUBE_LIVE")
	VIDEO_TYPE_SC       = MediaType("SOUNDCLOUD")
	VIDEO_TYPE_ODYSEE   = MediaType("ODYSEE")
	VIDEO_TYPE_MP4      = MediaType("MP4")
	VIDEO_TYPE_MP3      = MediaType("MP3")
	VIDEO_TYPE_PODCAST  = MediaType("PODCAST")
	VIDEO_TYPE_RG       = MediaType("RADIO_GARDEN")
	VIDEO_TYPE_PEERTUBE = MediaType("PEERTUBE")
)
