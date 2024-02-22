package utils

import (
	"encoding/json"
	"reflect"

	bolt "go.etcd.io/bbolt"
)

type Store[d any] struct {
	*bolt.DB
	Bucket []byte
}

func (store *Store[data]) Create() error {
	return store.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(store.Bucket)
		return err
	})
}

func (store *Store[data]) All() []data {
	values := []data{}
	store.View(func(tx *bolt.Tx) error {
		tx.Bucket(store.Bucket).ForEach(func(k, v []byte) error {
			var result data
			json.Unmarshal(v, &result)
			values = append(values, result)
			return nil
		})
		return nil
	})
	return values
}

func (store *Store[data]) Get(key string) (data, error) {
	var result data
	err := store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(store.Bucket).Get([]byte(key))
		return json.Unmarshal(b, &result)
	})
	return result, err
}

func (store *Store[data]) contains(obj data, field string, value any) bool {
	v := reflect.ValueOf(obj)

	fieldValue := v.Elem().FieldByName(field)
	if !fieldValue.IsValid() {
		return false // Field doesn't exist in the struct
	}
	return reflect.DeepEqual(fieldValue.Interface(), value)
}

func (store *Store[data]) Find(field string, value any) ([]data, error) {
	values := []data{}
	err := store.View(func(tx *bolt.Tx) error {
		tx.Bucket(store.Bucket).ForEach(func(k, v []byte) error {
			var result data
			err := json.Unmarshal(v, &result)
			if err != nil {
				return nil
			}
			if store.contains(result, field, value) {
				values = append(values, result)
			}
			return nil
		})
		return nil
	})
	return values, err
}

func (store *Store[data]) Delete(key string) error {
	return store.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(store.Bucket).Delete([]byte(key))
	})
}

func (store *Store[data]) Save(key string, value data) error {
	rawData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return store.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(store.Bucket).Put([]byte(key), rawData)
	})
}
