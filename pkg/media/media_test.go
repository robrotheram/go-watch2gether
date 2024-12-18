package media

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMP4Meida(t *testing.T) {
	client := MediaFactory.GetFactory("https://f000.backblazeb2.com/file/exceptionerror-io-public/movies/Bicentennial.Man.1999.mkv")
	tracks, err := client.GetMedia("https://f000.backblazeb2.com/file/exceptionerror-io-public/movies/Bicentennial.Man.1999.mkv", "")
	assert.Nil(t, err)
	assert.Equal(t, len(tracks), 1)
}

func TestPodcastMeida(t *testing.T) {
	url := "https://feeds.fireside.fm/selfhosted/rss"
	client := MediaFactory.GetFactory(url)
	tracks, err := client.GetMedia(url, "")
	assert.Nil(t, err)
	assert.Equal(t, len(tracks), 1)
}

func TestPeerTubeLive(t *testing.T) {
	url := "https://jupiter.tube/w/gYdezxbqGJMA3cfTTiK1cz"
	client := MediaFactory.GetFactory(url)
	tracks, err := client.GetMedia(url, "")
	assert.Nil(t, err)
	assert.Equal(t, len(tracks), 1)
}

func TestOdesey(t *testing.T) {
	url := "https://odysee.com/@veritasium:f/the-trillion-dollar-equation:3"
	client := MediaFactory.GetFactory(url)
	tracks, err := client.GetMedia(url, "")
	assert.Nil(t, err)
	assert.Equal(t, len(tracks), 1)
}

func TestYoutube(t *testing.T) {
	// url := "https://www.youtube.com/watch?v=kW8L7MVaaFE"
	url := "https://www.youtube.com/watch?v=WssuSLZJ9mY"
	client := MediaFactory.GetFactory(url)
	tracks, err := client.GetMedia(url, "")
	assert.Nil(t, err)
	assert.Equal(t, len(tracks), 1)

	fmt.Println(tracks[0].AudioUrl)
}
