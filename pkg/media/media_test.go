package media

import "testing"

func TestMeida(t *testing.T) {

	client := MediaFactory.GetFactory("https://f000.backblazeb2.com/file/exceptionerror-io-public/movies/Bicentennial.Man.1999.mkv")
	client.GetMedia("https://f000.backblazeb2.com/file/exceptionerror-io-public/movies/Bicentennial.Man.1999.mkv", "")
}
