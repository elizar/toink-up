package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type status struct {
	Status   string
	Time     int64
	Location string
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
		w.Write([]byte("404 - Not Found!"))
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

func getStatuses(courrier, trackingNumber string) (statuses []*status, err error) {
	switch courrier {
	case "phlpost":
		statuses, err = phlPost(trackingNumber)
	}

	return
}

func phlPost(tn string) (statuses []*status, err error) {
	const endpoint = "https://tnt.phlpost.gov.ph/"

	// Params
	p := make(url.Values)
	p.Add("TrackingNo", tn) // "CY023837389US"

	resp, err := http.PostForm(endpoint, p)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}

	// Loop through each row
	statuses = []*status{}
	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		// Exclude header
		if i == 0 {
			return
		}

		columns := []string{}
		row.Children().Each(func(i int, col *goquery.Selection) {
			columns = append(columns, strings.Trim(col.Text(), "  "))
		})

		// Parse time and subtract offset UTC+8 since PHLPOST
		// is using PH timezone
		t, _ := time.Parse("Jan 02 2006 3:04PM", columns[1])
		t = t.Add(-8 * time.Hour)

		statuses = append(statuses, &status{
			columns[0],
			t.UTC().Unix(),
			columns[2],
		})
	})

	return
}
