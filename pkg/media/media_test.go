package media

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMP4Meida(t *testing.T) {
	// client := MediaFactory.GetFactory("https://f000.backblazeb2.com/file/exceptionerror-io-public/movies/Bicentennial.Man.1999.mkv")
	tracks, err := MediaFactory.getMedia("https://f000.backblazeb2.com/file/exceptionerror-io-public/movies/Bicentennial.Man.1999.mkv", "")
	assert.Nil(t, err)
	assert.Equal(t, len(tracks), 1)
}

func TestPodcastMeida(t *testing.T) {
	url := "https://feeds.fireside.fm/selfhosted/rss"

	tracks, err := MediaFactory.getMedia(url, "")
	assert.Nil(t, err)
	assert.Equal(t, len(tracks), 1)
}

// func TestPeerTubeLive(t *testing.T) {
// 	url := "https://jupiter.tube/w/gYdezxbqGJMA3cfTTiK1cz"

// 	tracks, err := MediaFactory.getMedia(url, "")
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(tracks), 1)
// }

func TestOdesey(t *testing.T) {
	url := "https://odysee.com/@veritasium:f/the-trillion-dollar-equation:3"

	tracks, err := MediaFactory.getMedia(url, "")
	assert.Nil(t, err)
	assert.Equal(t, len(tracks), 1)
	assert.Nil(t, Refresh(tracks[0]))

}

func TestYoutube(t *testing.T) {
	// url := "https://www.youtube.com/watch?v=kW8L7MVaaFE"
	url := "https://www.youtube.com/watch?v=WssuSLZJ9mY"

	tracks, err := MediaFactory.getMedia(url, "")
	assert.Nil(t, err)
	assert.Equal(t, len(tracks), 1)

	fmt.Println(tracks[0].AudioUrl)
}
