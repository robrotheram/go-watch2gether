package main

import (
	"fmt"
	"testing"
	"watch2gether/pkg/media"

	"github.com/stretchr/testify/assert"
)

var username = "test-username"

func TestYoutubeGetAudioURL(t *testing.T) {
	url := "https://www.youtube.com/watch?v=l9nh1l8ZIJQ"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), media.MediaType("YOUTUBE"))
	assert.NotEmpty(m[0].AudioUrl)
	fmt.Println("URL:" + m[0].AudioUrl)
}

func TestYoutubeGetAudioOtherURL(t *testing.T) {
	url := "https://www.youtube.com/watch?v=k3WkJq478To"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), media.MediaType("YOUTUBE"))
	assert.NotEmpty(m[0].AudioUrl)
	fmt.Println("URL:" + m[0].AudioUrl)
}

func TestYoutubeGetLiveAudioURL(t *testing.T) {
	url := "https://www.youtube.com/watch?v=jfKfPfyJRdk"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), media.MediaType("YOUTUBE"))
	assert.NotEmpty(m[0].AudioUrl)
}

func TestSoudCloudGetAudioURL(t *testing.T) {
	url := "https://soundcloud.com/guy-j/guy_las-palapas"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), media.MediaType("SOUNDCLOUD"))
	assert.NotEmpty(m[0].AudioUrl)
}

func TestRadioGardenGetAudioURL(t *testing.T) {
	url := "https://radio.garden/listen/the-source/JL0Q8bRp"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), media.MediaType("RADIO_GARDEN"))
	assert.NotEmpty(m[0].AudioUrl)
}
func TestPodcastGetAudioURL(t *testing.T) {
	url := "https://feeds.fireside.fm/coder/rss"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), media.MediaType("PODCAST"))
	assert.NotEmpty(m[0].AudioUrl)
}

func TestMP3GetAudioURL(t *testing.T) {
	url := "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), media.MediaType("MP3"))
	assert.NotEmpty(m[0].AudioUrl)
}

func TestMP4GetAudioURL(t *testing.T) {
	url := "https://test-videos.co.uk/vids/bigbuckbunny/mp4/h264/1080/Big_Buck_Bunny_1080_10s_1MB.mp4"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), media.MediaType("MP4"))
	assert.NotEmpty(m[0].AudioUrl)
}

func TestPeerTubeAudioURL(t *testing.T) {
	url := "https://jupiter.tube/w/hNMXcK3L9XzQ5xxQj1YLXP"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), media.MediaType("PEERTUBE"))
	assert.NotEmpty(m[0].AudioUrl)
}

func TestODYSEEAudioURL(t *testing.T) {
	url := "https://odysee.com/@pianomusic:b/Christmas-Again:4"
	assert := assert.New(t)
	m, _ := media.MediaFactory.GetMedia(url, username)
	assert.Equal(len(m), 1)
	assert.Equal(media.MediaType("ODYSEE"), m[0].GetType())
	assert.NotEmpty(m[0].AudioUrl)
}
