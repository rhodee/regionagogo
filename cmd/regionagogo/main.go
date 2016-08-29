package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/akhenakh/regionagogo"
)

type Server struct {
	*regionagogo.GeoSearch
}

// countryHandler takes a lat & lng query params and return a JSON
// with the country of the coordinate
func (s *Server) countryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	slat := query.Get("lat")
	lat, err := strconv.ParseFloat(slat, 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	slng := query.Get("lng")
	lng, err := strconv.ParseFloat(slng, 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	region := s.StubbingQuery(lat, lng)
	w.Header().Set("Content-Type", "application/json")

	if region == nil {
		js, _ := json.Marshal(map[string]string{"name": "unknown"})
		w.Write(js)
		return
	}

	js, _ := json.Marshal(region.Data)
	w.Write(js)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbpath := flag.String("dbpath", "", "Database path")
	debug := flag.Bool("debug", false, "Enable debug")

	flag.Parse()

	gs, err := regionagogo.NewGeoSearch(*dbpath)
	if err != nil {
		log.Fatal(err)
	}
	gs.Debug = *debug

	err = gs.ImportGeoData()
	if err != nil {
		log.Fatal(err)
	}
	s := &Server{GeoSearch: gs}
	http.HandleFunc("/country", s.countryHandler)
	log.Println(http.ListenAndServe(":8082", nil))
}
