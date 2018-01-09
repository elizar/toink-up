package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/elizar/toink-up/parcel"
)

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
		if regexp.MustCompile("^\\/parcels").MatchString(r.URL.Path) {
			var err error
			var p interface{}

			code := http.StatusOK

			defer func() {
				// If has error :D
				if err != nil {
					p = struct {
						Code  int
						Error string
					}{code, err.Error()}
				}

				sb, _ := json.MarshalIndent(p, "", "  ")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
				w.Write(sb)
			}()

			// Excude the first item since it's empty
			segments := strings.Split(r.URL.Path, "/")[1:]

			if len(segments) != 3 {
				err = errors.New("Invalid request")
				code = http.StatusBadRequest
				return
			}

			// courier        -> segments[1]
			// trackingNumber -> segments[2]
			p = parcel.NewParcel(segments[1], segments[2])

			// Cast and fetch
			_, err = p.(*parcel.Parcel).Fetch()
			if err != nil {
				code = http.StatusBadRequest
				// not found
				rx := regexp.MustCompile(`(?i)not found`)
				if rx.MatchString(err.Error()) {
					code = http.StatusNotFound
				}
			}

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
