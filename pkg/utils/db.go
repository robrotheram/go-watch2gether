package utils

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// func CacheHashKey(tokenID string) string {
// 	return "app:auth:" + tokenID
// }

// func CacheHashField() string {
// 	return "token"
// }

func DBConnect(config Config) (*redis.Client, error) {
	fmt.Println(config.REDISURL)
	opt, err := redis.ParseURL(config.REDISURL)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(opt)
	return rdb, nil
}

func RedisConnect() {}

func RethinkDBConnect(config Config) (*rethinkdb.Session, error) {

	log.Infof("DB connection: %s Database: %s", config.RethinkDB, config.RethinkDatabase)

	session, err := rethinkdb.Connect(rethinkdb.ConnectOpts{
		Address:  config.RethinkDB, // endpoint without http
		Database: config.RethinkDatabase,
	})
	rethinkdb.DB(config.RethinkDatabase).TableCreate("user").Exec(session)
	rethinkdb.DB(config.RethinkDatabase).TableCreate("room").Exec(session)
	rethinkdb.DB(config.RethinkDatabase).TableCreate("hub").Exec(session)
	return session, err
}
