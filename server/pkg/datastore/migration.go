package datastore

import (
	"fmt"
	"log"
	"sort"

	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const MIGRATION_PREFIX = "migrations"

type migration interface {
	Migrate(data *Datastore) error
}

var MigrationFactory = make(map[string]migration)

type MigrationStore struct {
	session *rethinkdb.Session
}

func NewMigrationStore(session *rethinkdb.Session) *MigrationStore {
	rs := &MigrationStore{session: session}
	return rs
}

type MigrationVersion struct {
	Version string
}

func (udb *MigrationStore) Save(migration MigrationVersion) error {
	_, err := rethinkdb.Table(MIGRATION_PREFIX).Update(migration).RunWrite(udb.session)
	return err
}
func (udb *MigrationStore) Create(migration MigrationVersion) error {
	_, err := rethinkdb.Table(MIGRATION_PREFIX).Insert(migration).RunWrite(udb.session)
	return err
}

func (udb *MigrationStore) Get() (MigrationVersion, error) {
	version := MigrationVersion{}
	// Fetch all the items from the database
	res, err := rethinkdb.Table(MIGRATION_PREFIX).Run(udb.session)
	if err != nil {
		fmt.Println(err)
		return version, err
	}
	err = res.One(&version)

	if err != nil {
		fmt.Println(err)
		return version, err
	}
	return version, nil
}

func MigrationVersions() []string {
	keys := make([]string, len(MigrationFactory))
	i := 0
	for k := range MigrationFactory {
		keys[i] = k
		i++
	}
	return keys
}

func (data *Datastore) RunMigrations() {

	versions := MigrationVersions()
	sort.Strings(sort.StringSlice(versions))

	currentVersion, _ := data.Migrations.Get()
	if currentVersion.Version == "" {
		data.Migrations.Create(currentVersion)
	}
	pos := -1

	log.Println("Current Version: " + currentVersion.Version)
	for i, version := range versions {
		if currentVersion.Version == version {
			pos = i
		}
	}
	if pos+1 >= len(versions) {
		return
	}

	for _, version := range versions {
		log.Println("Running Migration for version: " + version)
		MigrationFactory[version].Migrate(data)
	}

	data.Migrations.Save(MigrationVersion{Version: VERSION})
}
