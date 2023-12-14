package stream

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"watch2gether/pkg/media"
)

type Segment struct {
	id        string
	index     float32
	videoFile string
	length    float32
}

func NewSegment(index int64, length float32, video media.Media) Segment {
	h := sha1.New()
	h.Write([]byte(video.Url))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return Segment{
		id:        fmt.Sprintf("%s-%d", sha1_hash, index),
		index:     float32(index),
		videoFile: video.AudioUrl,
		length:    length,
	}
}

func (s *Segment) ID() string {
	return s.id
}

func (s *Segment) startTime() float32 {
	return s.index * s.length
}

func (s *Segment) Args() []string {
	return []string{
		// Prevent encoding to run longer than 30 seonds
		"-timelimit", "45",

		// TODO: Some stuff to investigate
		// "-probesize", "524288",
		// "-fpsprobesize", "10",
		// "-analyzeduration", "2147483647",
		// "-hwaccel:0", "vda",

		// The start time
		// important: needs to be before -i to do input seeking
		"-ss", fmt.Sprintf("%v.00", s.startTime()),

		// The source file
		"-i", s.videoFile,

		// Put all streams to output
		// "-map", "0",

		// The duration
		"-t", fmt.Sprintf("%v.00", s.length),

		// TODO: Find out what it does
		//"-strict", "-2",

		// Synchronize audio
		"-async", "1",

		// 720p
		//"-vf", fmt.Sprintf("scale=-2:%v", res),

		// x264 video codec
		"-vcodec", "libx264",

		// x264 preset
		"-preset", "veryfast",

		// aac audio codec
		"-c:a", "aac",
		"-b:a", "128k",
		"-ac", "2",

		// TODO
		"-pix_fmt", "yuv420p",

		//"-r", "25", // fixed framerate

		//"-force_key_frames", "expr:gte(t,n_forced*5.000)",

		//"-force_key_frames", "00:00:00.00",
		//"-x264opts", "keyint=25:min-keyint=25:scenecut=-1",

		//"-f", "mpegts",

		"-f", "ssegment",
		"-segment_time", fmt.Sprintf("%v.00", s.length),
		"-initial_offset", fmt.Sprintf("%v.00", s.startTime()),
		"pipe:out%03d.ts",
	}
}
