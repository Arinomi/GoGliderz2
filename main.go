package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/marni/goigc"
	"net/http"
	"os"
	"time"
)

// global variable init
var startTime = time.Now().Truncate(time.Second) // starting time
var trackMAP = make(map[int]trackInfo)           // map of trackInfo structs
var ids = make([]int, 0)                         // array of ids

// struct init
type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

type trackInfo struct {
	Date     time.Time `json:"date"`
	Pilot    string    `json:"pilot"`
	Glider   string    `json:"glider"`
	GliderID string    `json:"glider_id"`
	Distance float64   `json:"distance"`
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

// creates a new track from the presented url and returns its ID
// returns 0 if the url was invalid
func newTrack(url string) int {
	newTrack, err := igc.ParseLocation(url)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	track := trackInfo{
		newTrack.Date,
		newTrack.Pilot,
		newTrack.GliderType,
		newTrack.GliderID,
		distance(newTrack)}

	trackMAP[len(trackMAP)+1] = track
	ids = append(ids, len(ids)+1)
	return len(ids)
}

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

	// routes init
	router.GET("/igcinfo/api", handlerAPI)
	router.GET("/igcinfo/api/igc", handlerIGC)
	router.GET("/igcinfo/api/igc/:id", handlerID)
	router.POST("/igcinfo/api/igc", handlerIGC)
	router.GET("/igcinfo/api/igc/:id/:field", handlerField)

	// server init
	http.ListenAndServe(getPort(), router)
}
