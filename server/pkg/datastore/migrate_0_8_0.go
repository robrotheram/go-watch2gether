package datastore

import (
	"watch2gether/pkg/media"

	log "github.com/sirupsen/logrus"
)

func init() {
	MigrationFactory["0.8.0"] = version_1{}
}

type version_1 struct{}

func (version_1) Migrate(data *Datastore) error {

	log.Info("Migrating Playlists to new version")
	playlist, _ := data.Playlist.GetAll()
	for _, playlist := range playlist {
		log.Info("Migrating Playlists:" + playlist.Name)
		for i := range playlist.Videos {
			v := &playlist.Videos[i]
			factory := media.MediaFactory.GetFactory(v.Url)
			v.Type = media.MediaType(factory.GetType())

			if v.Type == media.VIDEO_TYPE_YT {
				newMedia := factory.GetMedia(v.Url, playlist.Username)
				playlist.Videos[i] = newMedia[0]
			}
		}
		data.Playlist.Update(playlist)
	}
	return nil
}
