package utils

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	bolt "go.etcd.io/bbolt"
)

type testData struct {
	Name   string
	Age    int
	City   string
	Email  string
	Active bool
}

var (
	names  = []string{"Alice", "Bob", "Charlie", "David", "Emma", "Frank", "Grace", "Henry", "Ivy", "Jack"}
	cities = []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
	emails = []string{"gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "aol.com", "icloud.com"}
)

func generateRandomPerson() *testData {
	name := names[rand.Intn(len(names))]
	age := rand.Intn(50) + 18 // Random age between 18 and 67
	city := cities[rand.Intn(len(cities))]
	email := fmt.Sprintf("%s_%d@%s", name, rand.Intn(100), emails[rand.Intn(len(emails))])
	active := rand.Intn(2) == 0 // Random boolean

	return &testData{
		Name:   name,
		Age:    age,
		City:   city,
		Email:  email,
		Active: active,
	}
}

func TestStore(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test.db")
	if err != nil {
		t.Errorf("Error creating temp file: %v", err)
	}
	db, err := bolt.Open(tempFile.Name(), 0600, nil)
	if err != nil {
		t.Errorf("unable to create db: %v", err)
	}

	store := Store[*testData]{
		DB:     db,
		Bucket: []byte("test"),
	}
	if store.Create() != nil {
		t.Errorf("unable to create store: %v", err)
	}

	//saving
	if store.Save("test-1", generateRandomPerson()) != nil {
		t.Errorf("unable to save: %v", err)
	}
	if store.Save("test-2", generateRandomPerson()) != nil {
		t.Errorf("unable to save: %v", err)
	}
	if store.Save("test-3", generateRandomPerson()) != nil {
		t.Errorf("unable to save: %v", err)
	}
	if store.Save("test-4", generateRandomPerson()) != nil {
		t.Errorf("unable to save: %v", err)
	}
	if store.Save("test-5", generateRandomPerson()) != nil {
		t.Errorf("unable to save: %v", err)
	}

	if person, err := store.Get("test-5"); err != nil || len(person.Name) == 0 {
		t.Errorf("unable to get: %v", err)
	}

	if size := len(store.All()); size != 5 {
		t.Errorf("unable to list expected 5 got %d", size)
	}

	//update
	person, _ := store.Get("test-4")
	person.Name = "Martin"
	store.Save("test-4", person)

	if person, err := store.Get("test-4"); err != nil || person.Name != "Martin" {
		t.Errorf("unable to Update Name should be Martin was: %s", person.Name)
	}

	//delete
	if err := store.Delete("test-3"); err != nil {
		t.Errorf("unable to delete: %v", err)
	}
	if size := len(store.All()); size != 4 {
		t.Errorf("unable to delete expected 4 got %d", size)
	}

	//find
	if list, _ := store.Find("Name", "Martin"); len(list) != 1 {
		t.Errorf("unable to Find expected 1 got %d", len(list))
	}

}
