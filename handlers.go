package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/marni/goigc"
	"net/http"
	"time"
)

func handlerRedir(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	newPath := r.URL.Path + "/api"
	http.Redirect(w, r, newPath, http.StatusPermanentRedirect)
}

// handler for igcinfo/api
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

func handlerTrack(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	switch r.Method {
	case "GET":
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	case "POST":
		input := make(map[string]interface{})
		_ = json.NewDecoder(r.Body).Decode(&input)

		newID := "id" + string(db.Count()+1)

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
