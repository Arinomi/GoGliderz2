package main

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/julienschmidt/httprouter"
	"github.com/marni/goigc"
	"net/http"
	"strconv"
	"time"
)

// redirection handler for paragliding/ to paragliding/api
func handlerRedir(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	newPath := r.URL.Path + "/api"
	http.Redirect(w, r, newPath, http.StatusPermanentRedirect)
}

// handler for paragliding/api
func handlerAPI(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	info := apiInfo{uptime(), "Service for Paragliding tracks.", "v1"}

	jsresp, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsresp)
}

// handler for paragliding/api/track
func handlerTrack(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	switch r.Method {
	case "GET":
		tracks := db.GetAll()
		var trackIDs []string
		for _, element := range tracks {
			trackIDs = append(trackIDs, element.TrackID)
		}

		jsresp, err := json.Marshal(trackIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(jsresp)

	case "POST":
		input := make(map[string]interface{})
		_ = json.NewDecoder(r.Body).Decode(&input)

		newID := "id" + strconv.Itoa(db.Count()+1)

		newTrack, err := igc.ParseLocation(input["url"].(string))
		if err != nil {
			http.Error(w, "Not able to process the URL", http.StatusBadRequest)
			return
		}

		newInsert := Track{
			time.Now(),
			newID,
			newTrack.Date,
			newTrack.Pilot,
			newTrack.GliderType,
			newTrack.GliderID,
			distance(newTrack),
			input["url"].(string)}

		err = db.Add(newInsert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsresp, err := json.Marshal(map[string]string{"id": newID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsresp)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

// handler for paragliding/api/track/<id>
func handlerTrackID(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	id := ps[0].Value
	fmt.Println(ps[0].Value)

	track, ok := db.GetSelect(id, sel("H_date", "pilot", "glider", "glider_id", "track_length", "track_src_url"))
	if ok {
		if len(track) != 0 {
			jsresp, err := json.Marshal(track)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Write(jsresp)
		} else {
			http.Error(w, "No track found with that ID", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Retrieving the track returned an error", http.StatusInternalServerError)
	}
}

// handler for paragliding/api/track/<id>/field
func handlerTrackField(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	id := ps[0].Value
	field := ps[1].Value

	track, ok := db.GetSelect(id, sel(field))
	if ok {
		if len(track) != 0 {
			jsresp, err := json.Marshal(track)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Write(jsresp)
		} else {
			http.Error(w, "No track found with that ID", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Retrieving the track returned an error", http.StatusInternalServerError)
	}
}

// handler for paragliding/api/ticker
func handlerTicker(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	ticker, ok := returnTicker()
	if !ok {
		http.Error(w, "No tracks found", http.StatusNotFound)
	} else {
		jsresp, err := json.Marshal(ticker)
		if err != nil {
			http.Error(w, "Error parsing the Ticker", http.StatusInternalServerError)
		}
		w.Write(jsresp)
	}
}

// handler for paragliding/api/ticker/<timestamp> and /latest
func handlerTickerTimestamps(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	input := ps[0].Value

	if input == "latest" {
		timestamps := returnTimestamps(time.Time{})
		latest := timestamps.latest

		if latest.IsZero() {
			http.Error(w, "No records found", http.StatusBadRequest)
		} else {
			query := bson.M{"timestamp": latest}
			track, ok := db.Get(query)
			if ok {
				jsresp, err := json.Marshal(track.Timestamp)
				if err != nil {
					http.Error(w, "Error parsing timestamp", http.StatusInternalServerError)
				}
				w.Write(jsresp)
			}
		}
	} else {
		timeInput, err := time.Parse(time.RFC3339, input) // Check if the timestamp provided is a valid time
		if err != nil {
			http.Error(w, "Invalid time format, use DD.MM.YYYY HH:MM:SS", http.StatusBadRequest)
			return
		}

		ticker, ok := returnTickerTimestamp(timeInput)
		if !ok {
			http.Error(w, "No tracks found", http.StatusNotFound)
		} else {
			jsresp, err := json.Marshal(ticker)
			if err != nil {
				http.Error(w, "Error parsing the Ticker", http.StatusInternalServerError)
			}
			w.Write(jsresp)
		}
	}
}
