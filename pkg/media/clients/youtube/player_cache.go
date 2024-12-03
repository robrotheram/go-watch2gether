package youtube

import (
	"time"
)

const defaultCacheExpiration = time.Minute * time.Duration(5)

type playerCacheRow struct {
	key       string
	expiredAt time.Time
	config    playerConfig
}

type playerCache []playerCacheRow

// Get : get cache  when it has same video id and not expired
func (s playerCache) Get(key string) playerConfig {
	return s.GetCacheBefore(key, time.Now())
}

// GetCacheBefore : can pass time for testing
func (s playerCache) GetCacheBefore(key string, time time.Time) playerConfig {
	for _, cache := range s {
		if key == cache.key && cache.expiredAt.After(time) {
			return cache.config
		}
	}
	return nil
}

// Set : set cache with default expiration
func (s *playerCache) Set(key string, operations playerConfig) {
	s.setWithExpiredTime(key, operations, time.Now().Add(defaultCacheExpiration))
}

func (s *playerCache) setWithExpiredTime(key string, config playerConfig, time time.Time) {
	haskey := false
	for _, cache := range *s {
		if key == cache.key && cache.expiredAt.After(time) {
			haskey = true
		}
	}
	if !haskey {
		*s = append(*s, playerCacheRow{
			key:       key,
			config:    config,
			expiredAt: time,
		})
	}
}
