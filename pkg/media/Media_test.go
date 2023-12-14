package media

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var username = "test-username"

func TestYoutubeGetAudioURL(t *testing.T) {
	url := "https://www.youtube.com/watch?v=RlGOKUFZqGo"
	assert := assert.New(t)
	m, _ := MediaFactory.GetMedia(url, username)
	RefreshAudioURL(&m[0])
	fmt.Println("URL:" + m[0].AudioUrl)

	assert.Equal(len(m), 1)
	assert.Equal(m[0].GetType(), MediaType("YOUTUBE"))
	assert.NotEmpty(m[0].AudioUrl)

}
