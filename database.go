package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// StudentsMongoDB stores the details of the DB connection.
type tracksMongoDB struct {
	DatabaseURL          string
	DatabaseName         string
	TracksCollectionName string
}

/*
Init initializes the mongo storage.
*/
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

/*
Add adds new students to the storage.
*/
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

/*
Count returns the current count of the students in in-memory storage.
*/
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

/*
Get returns a student with a given ID or empty student struct.
*/
func (db *tracksMongoDB) Get(keyID string) (Track, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	student := Track{}
	ok := true

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).Find(bson.M{"track_id": keyID}).One(&student)
	if err != nil {
		ok = false
	}

	return student, ok
}

/*
GetAll returns a slice with all the students.
*/
func (db *tracksMongoDB) GetAll() []Track {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var all []Track

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).Find(bson.M{}).All(&all)
	if err != nil {
		return []Track{}
	}

	return all
}
