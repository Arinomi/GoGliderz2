package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

// handler for igcinfo/api
func handlerAPI(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	info := apiInfo{uptime(), "Service for IGC tracks.", "v1"}

	jsresp, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsresp)
}

// handler for igcinfo/api/igc
func handlerIGC(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	switch r.Method {
	case "GET":
		jsresp, err := json.Marshal(ids)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsresp)
	case "POST":
		input := make(map[string]interface{})
		_ = json.NewDecoder(r.Body).Decode(&input)

		newID := newTrack(input["url"].(string))
		if newID == 0 {
			http.Error(w, "Not able to process the URL", http.StatusBadRequest)
			return
		}

		jsresp, err := json.Marshal(map[string]int{"id": newID})
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

// handler for igcinfo/api/igc/<id>
func handlerID(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	if len(trackMAP) > 0 && len(ids) > 0 {
		http.Header.Add(w.Header(), "content-type", "application/json")

		id, err := strconv.Atoi(ps[0].Value)
		if err != nil {
			http.Error(w, "Please provide a valid ID.", http.StatusBadRequest)
			return
		} else {
			_, ok := trackMAP[id]
			if ok {
				fmt.Println(id)
				jsresp, err := json.Marshal(trackMAP[id])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(jsresp)
			} else {
				http.Error(w, "Given ID not found", http.StatusNotFound)
			}
		}
	} else {
		http.Error(w, "No files found", http.StatusNotFound)
	}
}

// handler for igcinfo/api/<id>/<field>
func handlerField(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	if len(trackMAP) > 0 && len(ids) > 0 {

		id, err := strconv.Atoi(ps[0].Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		field := ps[1].Value

		switch field {
		case "pilot":
			w.Write([]byte(trackMAP[id].Pilot))

		case "glider":
			w.Write([]byte(trackMAP[id].Glider))

		case "glider_id":
			w.Write([]byte(trackMAP[id].GliderID))

		case "calculated total track length":
			distString := strconv.FormatFloat(trackMAP[id].Distance, 'f', -1, 64)
			w.Write([]byte(distString))

		case "H_date":
			w.Write([]byte(trackMAP[id].Date.Format(time.RFC3339)))

		default:
			http.Error(w, "Not a valid field", http.StatusBadRequest)
		}

	} else {
		http.Error(w, "No files found", http.StatusNotFound)
	}
}
