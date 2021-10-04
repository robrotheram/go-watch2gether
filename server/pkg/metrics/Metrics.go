package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalRoomsCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "watch2gether_total_rooms",
		Help: "Number of rooms",
	})
	TotalUsersCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "watch2gether_total_users",
		Help: "Number of users",
	})
	ActiveRoomsCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "watch2gether_active_rooms",
		Help: "Number of pings",
	})
	AudioCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "watch2gether_concurrent_audio_streams",
		Help: "Number of pings",
	})
	BotCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "watch2gether_active_bots",
		Help: "Number of pings",
	})
	ActiveUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "watch2gether_active_users",
		Help: "Number of pings",
	})
	OpsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "watch2gether_processed_ops_total",
		Help: "The total number of processed events",
	})
)
