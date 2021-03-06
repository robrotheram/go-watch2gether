package datastore

import (
	"watch2gether/pkg/utils"

	log "github.com/sirupsen/logrus"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func RedisConnect() {}

func createSession(config utils.Config) (*rethinkdb.Session, error) {

	log.Infof("DB connection: %s Database: %s", config.RethinkURL, config.RethinkDatabase)

	session, err := rethinkdb.Connect(rethinkdb.ConnectOpts{
		Address:  config.RethinkURL, // endpoint without http
		Database: config.RethinkDatabase,
	})
	rethinkdb.DBCreate(config.RethinkDatabase).Exec(session)
	return session, err
}

func createTable(session *rethinkdb.Session, config utils.Config, table string) error {
	return rethinkdb.DB(config.RethinkDatabase).TableCreate(table).Exec(session)
}
