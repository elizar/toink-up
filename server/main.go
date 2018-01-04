package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type parcel struct {
	Status    string
	Time      int64
	Location  string
	Delivered bool
}

func main() {
	// Init server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Index page
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./public/index.html")
			return
		}

		// Serve static files
		if regexp.MustCompile("\\..{2,}$").MatchString(r.URL.Path) {
			http.ServeFile(w, r, fmt.Sprintf("%s/%s", "./public/", r.URL.Path))
			return
		}

		// Tracking and shit
		if regexp.MustCompile("^\\/parcels").MatchString(r.URL.Path) && r.Method == http.MethodPost {
			segments := strings.Split(r.URL.Path, "/")[1:] // Ignore the first item which is an empty string

			if len(segments) != 3 {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			// courier => segments[1]
			// tracking_number => segments[2]
			statuses, err := getStatuses(segments[1], segments[2])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			sb, _ := json.MarshalIndent(statuses, "", " ")
			w.Header().Set("Content-Type", "application/json")
			w.Write(sb)

			return
		}

		// Not found
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html")
		f, _ := os.Open("./public/404.html")
		b, _ := ioutil.ReadAll(f)
		w.Write(b)
	})

	// Configure default port
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	// Listen and server mo'fucker!
	log.Println("[ Server ] - up and running on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func getStatuses(courrier, trackingNumber string) (statuses []*parcel, err error) {
	switch courrier {
	case "phlpost":
		statuses, err = phlPost(trackingNumber)
	}

	return
}
