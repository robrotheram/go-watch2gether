package datastore

import (
	"time"
	"watch2gether/pkg/metrics"

	log "github.com/sirupsen/logrus"
)

type MetricCollections struct {
	datastore Datastore
	ticker    time.Ticker
	done      chan bool
}

func NewMetricCollection(ds Datastore) *MetricCollections {
	return &MetricCollections{
		datastore: ds,
		ticker:    *time.NewTicker(time.Second * 2),
		done:      make(chan bool),
	}
}

func (mc *MetricCollections) Start() {
	go func() {
		for {
			select {
			case <-mc.done:
				log.Debug("DONE")
				return
			case <-mc.ticker.C:
				mc.ClollectMetrics()
			}
		}
	}()
}
func (mc *MetricCollections) Stop() {
	mc.ticker.Stop()
	mc.done <- true
}

func (mc *MetricCollections) ClollectMetrics() {

	if data, err := mc.datastore.Rooms.GetAll(); err == nil {
		metrics.TotalRoomsCounter.Set(float64(len(data)))
	}
	if data, err := mc.datastore.Users.GetAll(); err == nil {
		metrics.TotalUsersCounter.Set(float64(len(data)))
	}

	rooms := mc.datastore.Hub.Rooms
	metrics.ActiveRoomsCounter.Set(float64(len(rooms)))
	activeUsers := float64(0)
	acticeAudio := float64(0)
	activeBots := float64(0)

	for _, room := range rooms {
		activeUsers = activeUsers + float64(len(room.Clients))
		if room.Bot != nil {
			activeBots = activeBots + 1.0
			if room.Bot.Audio != nil {
				if room.Bot.Audio.Playing {
					acticeAudio = acticeAudio + 1.0
				}
			}
		}
	}
	metrics.ActiveUsers.Set(activeUsers)
	metrics.AudioCounter.Set(acticeAudio)
	metrics.BotCounter.Set(activeBots)

}
