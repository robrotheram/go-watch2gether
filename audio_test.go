package main

import (
	"testing"
	"watch2gether/pkg/audioBot"
)

func TestHello(t *testing.T) {
	a := audioBot.Audio{Url: "https://filesamples.com/samples/video/mp4/sample_1280x720_surfing_with_audio.mp4"}
	duration, err := a.GetDuration()
	if err != nil {
		t.Errorf("Error Getting Data %w", err)
	}
	got := duration.Seconds()
	want := float64(183.147000)

	if got != want {
		t.Errorf("got %f want %f", got, want)
	}
}
