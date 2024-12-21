package history

import (
	"w2g/pkg/media"
	"w2g/pkg/utils"

	bolt "go.etcd.io/bbolt"
)

type History struct {
	Videos []*media.Media `json:"videos"`
}

type HistoryStore struct {
	*utils.Store[*History]
}

func NewHistoryStore(db *bolt.DB) *HistoryStore {
	store := &utils.Store[*History]{
		DB:     db,
		Bucket: []byte("history"),
	}
	store.Create()

	return &HistoryStore{
		Store: store,
	}
}

func (store *HistoryStore) getHisory(channel string) (*History, error) {
	return store.Get(channel)
}

func (store *HistoryStore) GetHisory(channel string) ([]*media.Media, error) {
	hisory, err := store.getHisory(channel)
	if err != nil {
		return []*media.Media{}, err
	}
	return hisory.Videos, nil
}

func (store *HistoryStore) AddTracks(channel string, tracks []*media.Media) error {
	history, err := store.getHisory(channel)
	if err != nil {
		history = &History{
			Videos: tracks,
		}
		return store.Save(channel, history)
	}

	history.Videos = removeDuplicates(append(tracks, history.Videos...))
	return store.Save(channel, history)
}

func (store *HistoryStore) AddTrack(channel string, tracks *media.Media) error {
	return store.AddTracks(channel, []*media.Media{tracks})
}

func removeDuplicates(scripts []*media.Media) []*media.Media {
	seen := make(map[string]bool)
	var result []*media.Media
	for _, script := range scripts {
		if !seen[script.Url] {
			seen[script.Url] = true
			result = append(result, script)
		}
	}
	return result
}
