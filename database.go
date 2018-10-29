package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// stores the details of the DB connection.
type tracksMongoDB struct {
	DatabaseURL          string
	DatabaseName         string
	TracksCollectionName string
}

// Initializes the connection
func (db *tracksMongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"track_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// add new track to the database
func (db *tracksMongoDB) Add(t Track) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).Insert(t)
	if err != nil {
		fmt.Printf("Error in Add(): %v", err.Error())
		return err
	}

	return nil
}

// return number of tracks in the database
func (db *tracksMongoDB) Count() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.TracksCollectionName).Count()
	if err != nil {
		fmt.Printf("Error in Count(): %v", err.Error())
		return -1
	}

	return count
}

// return a single track and a bool, requires an ID and a bson.M query to select which fields to return
func (db *tracksMongoDB) GetSelect(keyID string, fields bson.M) (map[string]interface{}, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	collection := session.DB(db.DatabaseName).C(db.TracksCollectionName)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true

	result := make(map[string]interface{})

	err = collection.Find(bson.M{"track_id": keyID}).Select(fields).One(&result)
	if err != nil {
		ok = false
	}

	return result, ok
}

// return a single track and a bool, requires a bson.M query to choose which fields to search by
func (db *tracksMongoDB) Get(query bson.M) (Track, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	collection := session.DB(db.DatabaseName).C(db.TracksCollectionName)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true

	var result Track

	err = collection.Find(query).One(&result)
	if err != nil {
		ok = false
	}

	return result, ok
}

// returns a slice with all the tracks
func (db *tracksMongoDB) GetAll() []Track {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var all []Track

	err = session.
		DB(db.DatabaseName).
		C(db.TracksCollectionName).
		Find(bson.M{}).
		All(&all)
	if err != nil {
		return []Track{}
	}

	return all
}

// returns a slice with all the tracks, sorted by the field given
func (db *tracksMongoDB) GetAllSorted(field string) []Track {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var all []Track

	err = session.
		DB(db.DatabaseName).
		C(db.TracksCollectionName).
		Find(bson.M{}).
		Sort(field).
		All(&all)
	if err != nil {
		return []Track{}
	}

	return all
}
