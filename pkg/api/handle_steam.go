package api

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os/exec"
	"strconv"
	"syscall"
	"watch2gether/pkg/api/stream"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const PlaylistContentType = "application/vnd.apple.mpegurl"
const CaptionContentType = "text/vtt"
const hlsSegmentLength = 20.0 // Seconds\

func (api *API) handleStream(c echo.Context) error {
	id := c.Param("id")

	controller, _ := api.GetController(id)
	state, err := controller.GetState()
	if err != nil {
		return err
	}
	state.Current.Refresh()
	controller.Save(state)
	duration := state.Current.Duration.Seconds()

	c.Response().Header().Set("Content-Type", PlaylistContentType)
	urlTemplate := fmt.Sprintf("%v://%v/api/channel/%s/stream/{{.Segment}}?id=%s", "http", c.Request().Host, id, state.Current.ID)
	t := template.Must(template.New("urlTemplate").Parse(urlTemplate))

	getUrl := func(segmentIndex int) string {
		buf := new(bytes.Buffer)
		t.Execute(buf, struct {
			Resolution int64
			Segment    int
		}{
			720,
			segmentIndex,
		})
		return buf.String()
	}

	w := c.Response().Writer

	fmt.Fprint(w, "#EXTM3U\n")
	fmt.Fprint(w, "#EXT-X-VERSION:3\n")
	fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprint(w, "#EXT-X-ALLOW-CACHE:YES\n")
	fmt.Fprint(w, "#EXT-X-TARGETDURATION:"+fmt.Sprintf("%.f", hlsSegmentLength)+"\n")
	fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	leftover := duration
	segmentIndex := 0

	for leftover > 0 {
		if leftover > hlsSegmentLength {
			fmt.Fprintf(w, "#EXTINF: %f,\n", hlsSegmentLength)
		} else {
			fmt.Fprintf(w, "#EXTINF: %f,\n", leftover)
		}
		fmt.Fprintf(w, getUrl(segmentIndex)+"\n")
		segmentIndex++
		leftover = leftover - hlsSegmentLength
	}
	fmt.Fprint(w, "#EXT-X-ENDLIST\n")
	return nil
}

func (api *API) handleSegment(c echo.Context) error {
	index, _ := strconv.ParseInt(c.Param("segment"), 0, 64)
	id := c.Param("id")
	controller, _ := api.GetController(id)
	state, err := controller.GetState()
	if err != nil {
		return err
	}
	segment := stream.NewSegment(index, hlsSegmentLength, state.Current)
	encoder := NewEncoder()
	return encoder.Serve(segment, c.Response().Writer)
}

func NewEncoder() *Encoder {
	return &Encoder{
		cache: stream.NewCache(),
	}
}

type Encoder struct {
	cache *stream.Cache
}

func (e *Encoder) tryServeFromCache(segment stream.Segment, w io.Writer) (bool, error) {
	data, err := e.cache.Get(context.Background(), segment.ID())
	if err != nil {
		return false, err
	}
	if data == nil {
		e.cache.Delete(segment.ID())
		return false, fmt.Errorf("no data from cache")
	}
	if _, err = io.Copy(w, bytes.NewReader(data)); err != nil {
		return true, err
	}
	return true, nil
}

func (e *Encoder) Serve(segment stream.Segment, w io.Writer) error {
	if served, err := e.tryServeFromCache(segment, w); served || err != nil {
		if served {
			log.Debugf("Served request %v from cache", segment.ID())
			return nil
		}
	}
	cw := new(bytes.Buffer)
	mw := io.MultiWriter(cw, w)
	cmd := exec.Command("ffmpeg", segment.Args()...)
	err := e.ExecAndWriteStdout(cmd, mw)
	if err != nil {
		log.Infof("FFMPEG: %v", err)
		return err
	}

	return e.cache.Set(context.Background(), segment.ID(), cw.Bytes())
}

func (e *Encoder) ExecAndWriteStdout(cmd *exec.Cmd, w io.Writer) error {
	stdout, err := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error opening stdout of command: %v", err)
	}
	defer stdout.Close()
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	if _, err := io.Copy(w, stdout); err != nil {
		// Ask the process to exit
		cmd.Process.Signal(syscall.SIGKILL)
		cmd.Process.Wait()
		return fmt.Errorf("error copying stdout to buffer: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Infof("FFMPEG ERR: %s", scanner.Text())
		}
		// log.Infof("FFMPEG ERR: %s", cmdDebug("ffmpeg", cmd.Args))
		return fmt.Errorf("command failed %v", err)
	}

	return nil
}

func cmdDebug(cmd string, args []string) string {
	for _, arg := range args {
		cmd += " " + arg
	}
	return cmd
}
