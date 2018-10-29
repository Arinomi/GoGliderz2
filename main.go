package main

import (
	"github.com/globalsign/mgo/bson"
	"github.com/julienschmidt/httprouter"
	"github.com/marni/goigc"
	"net/http"
	"os"
	"time"
)

// global variable init
var startTime = time.Now().Truncate(time.Second) // starting time

// database init
var db = tracksMongoDB{"mongodb://user:password123@ds141783.mlab.com:41783/tracks", "tracks", "tracks"}

// struct init

// apiInfo is written to http.ResponseWriter
type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// Track is written to http.ResponseWriter
type Track struct {
	Timestamp time.Time `bson:"timestamp"`
	TrackID   string    `bson:"track_id"`
	Date      time.Time `bson:"H_date"`
	Pilot     string    `bson:"pilot"`
	Glider    string    `bson:"glider"`
	GliderID  string    `bson:"glider_id"`
	Distance  float64   `bson:"track_length"`
	SrcURL    string    `bson:"track_src_url"`
}

// returning uptime as a string in ISO 8601/RFC3339 format
func uptime() string {
	now := time.Now().Truncate(time.Second)

	// formatting using the example string from the Wikipedia page
	now.Format("P3Y6M4DT12H30M5S")
	startTime.Format("P3Y6M4DT12H30M5S")

	return now.Sub(startTime).String()
}

// calculating the total distance using the track's points
func distance(p igc.Track) float64 {
	d := 0.0

	for i := 0; i < len(p.Points)-1; i++ {
		d += p.Points[i].Distance(p.Points[i+1])
	}

	return d
}

// generating a bson.M selection
func sel(q ...string) (r bson.M) {
	r = make(bson.M, len(q))
	r["_id"] = 0
	for _, s := range q {
		r[s] = 1
	}
	return
}

// getting the appropriate port for Heroku
func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":80"
}

// main function
func main() {
	// router init
	router := httprouter.New()

	db.Init()

	// routes init
	router.GET("/paragliding", handlerRedir)
	router.GET("/paragliding/api", handlerAPI)
	router.GET("/paragliding/api/track", handlerTrack)
	router.POST("/paragliding/api/track", handlerTrack)
	router.GET("/paragliding/api/track/:id", handlerTrackID)
	router.GET("/paragliding/api/track/:id/:field", handlerTrackField)
	router.GET("/paragliding/api/ticker/", handlerTicker)
	router.GET("/paragliding/api/ticker/:time/", handlerTickerTimestamps)

	// server init
	http.ListenAndServe(getPort(), router)
}
