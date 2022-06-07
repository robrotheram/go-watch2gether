package main

import (
	"testing"
	"watch2gether/pkg/media"

	"github.com/stretchr/testify/assert"
)

func TestGetMediaTypes(t *testing.T) {
	assert := assert.New(t)
	types := media.MediaFactory.GetTypes()
	assert.ElementsMatch(types, []string{"MP3", "MP4", "PODCAST", "RADIO_GARDEN", "SOUNDCLOUD", "YOUTUBE"})
}

func TestYoutubeTypes(t *testing.T) {
	url := "https://www.youtube.com/watch?v=l9nh1l8ZIJQ&t=2968s"
	assert := assert.New(t)
	factory := media.MediaFactory.GetFactory(url)
	assert.Equal(factory.GetType(), "YOUTUBE")
}
func TestSoudCloudType(t *testing.T) {
	url := "https://soundcloud.com/guy-j/guy_las-palapas"
	assert := assert.New(t)
	factory := media.MediaFactory.GetFactory(url)
	assert.Equal(factory.GetType(), "SOUNDCLOUD")
}

func TestRadioGardenType(t *testing.T) {
	url := "https://radio.garden/listen/the-source/JL0Q8bRp"
	assert := assert.New(t)
	factory := media.MediaFactory.GetFactory(url)
	assert.Equal(factory.GetType(), "RADIO_GARDEN")
}
func TestPodcastType(t *testing.T) {
	url := "https://feeds.fireside.fm/coder/rss"
	assert := assert.New(t)
	factory := media.MediaFactory.GetFactory(url)
	assert.Equal(factory.GetType(), "PODCAST")
}

func TestMP3Type(t *testing.T) {
	url := "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3"
	assert := assert.New(t)
	factory := media.MediaFactory.GetFactory(url)
	assert.Equal(factory.GetType(), "MP3")
}

func TestMP4Type(t *testing.T) {
	url := "https://test-videos.co.uk/vids/bigbuckbunny/mp4/h264/1080/Big_Buck_Bunny_1080_10s_1MB.mp4"
	assert := assert.New(t)
	factory := media.MediaFactory.GetFactory(url)
	assert.Equal(factory.GetType(), "MP4")
}

func TestCRAZY(t *testing.T) {
	url := "https://vs-cmaf-pushb-uk.live.cf.md.bbci.co.uk/x=3/i=urn:bbc:pips:service:bbc_one_west_midlands/pc_hd_abr_v2.mpd"
	assert := assert.New(t)
	factory := media.MediaFactory.GetFactory(url)
	assert.Equal(factory.GetType(), "MP4")
}

func TestJupiter(t *testing.T) {
	url := "https://jupiter.tube/w/sroVsxapHuMc4VjpgTb1em"
	assert := assert.New(t)
	factory := media.MediaFactory.GetFactory(url)
	assert.Equal(factory.GetType(), "PEERTUBE")
}

func TestJupiterM3u8(t *testing.T) {
	url := "https://jupiter.tube/static/streaming-playlists/hls/d619534d-4820-4ffa-a760-25e88d724ae6/master.m3u8"
	assert := assert.New(t)
	factory := media.MediaFactory.GetFactory(url)
	assert.Equal(factory.GetType(), "MP4")
}
