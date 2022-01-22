package datastore

import (
	"watch2gether/pkg/media"

	log "github.com/sirupsen/logrus"
)

func init() {
	MigrationFactory["0.8.0"] = verion_1{}
}

type verion_1 struct{}

func (verion_1) Migrate(data *Datastore) error {
	client := media.GetDownloader()
	log.Info("Migreating Playlists to new version")
	playlist, _ := data.Playlist.GetAll()
	for _, playist := range playlist {
		log.Info("Migrating Playlists:" + playist.Name)
		for i := range playist.Videos {
			v := &playist.Videos[i]
			v.Type, _ = media.TypeFromUrl(v.Url)
			if v.Type == media.VIDEO_TYPE_YT {
				ytVideo, err := client.GetVideo(v.Url)
				if err == nil {
					v.Update(ytVideo)
				}
			}
		}
		data.Playlist.Update(playist)
	}
	return nil
}
