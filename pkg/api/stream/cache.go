package stream

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/djherbis/atime"
	log "github.com/sirupsen/logrus"
)

type Cache struct {
	path string
}

func NewCache() *Cache {
	cache := &Cache{os.TempDir()}
	cache.Prune()
	return cache
}

func (d *Cache) Prune() error {
	currentTime := time.Now()
	cutoffTime := currentTime.Add(-time.Minute * 10) //10mins cache

	files, err := os.ReadDir(d.path)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(d.path, file.Name())
		at, err := atime.Stat(path)
		if err != nil {
			continue
		}
		diff := at.Before(cutoffTime)
		if diff {
			d.Delete(file.Name())
		}
	}
	return nil
}

func (d *Cache) getCachePath(key string) string {
	return filepath.Join(d.path, key)
}

func (d *Cache) Delete(key string) error {
	log.Debug("Removing file")
	path := filepath.Join(d.path, key)
	return os.Remove(path)
}

func (d *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	file, err := os.Open(d.getCachePath(key))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()
	b := new(bytes.Buffer)
	if _, err = io.Copy(b, file); err != nil {
		log.Errorf("Error copying file to cache value: %v", err)
		return nil, err
	}
	return b.Bytes(), nil
}

func (d *Cache) Set(ctx context.Context, key string, value []byte) error {
	log.Debugf("Setting cache item %v")
	if err := os.MkdirAll(d.path, 0777); err != nil {
		log.Errorf("Could not create cache dir %v: %v", d.path, err)
		return err
	}
	cacheTmpFile, err := os.CreateTemp(d.path, key+".*")
	if err != nil {
		log.Errorf("Could not create cache file %v: %v", cacheTmpFile, err)
		return err
	}
	if _, err := io.Copy(cacheTmpFile, bytes.NewReader(value)); err != nil {
		log.Errorf("Could not write cache file %v: %v", cacheTmpFile, err)
		cacheTmpFile.Close()
		os.Remove(cacheTmpFile.Name())
		return err
	}
	if err = cacheTmpFile.Close(); err != nil {
		log.Errorf("Could not close cache file %v: %v", cacheTmpFile, err)
		os.Remove(cacheTmpFile.Name())
		return err
	}
	if err = os.Rename(cacheTmpFile.Name(), d.getCachePath(key)); err != nil {
		log.Errorf("Could not move cache file %v: %v", cacheTmpFile, err)
		os.Remove(cacheTmpFile.Name())
		return err
	}
	return nil
}
