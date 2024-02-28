package utils

import (
	_ "embed"
	"runtime/debug"
	"time"
)

//go:generate sh -c "git describe --tags --abbrev=0 > version.txt"
//go:embed version.txt
var Version string
var LastCommit time.Time
var Revision string

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, kv := range info.Settings {
		if kv.Value == "" {
			continue
		}
		switch kv.Key {
		case "vcs.revision":
			Revision = kv.Value
		case "vcs.time":
			LastCommit, _ = time.Parse(time.RFC3339, kv.Value)
		}
	}
}
