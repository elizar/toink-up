package parcel

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// NewTracker ...
func NewParcel(courier, tn string) *Parcel {
	return &Parcel{
		Courier:        courier,
		TrackingNumber: tn,
	}
}

// Fetch retrives the tracking history of a given parcel
func (p *Parcel) Fetch() (total int, err error) {
	switch p.Courier {
	case PHLPOST:
		var statuses []*Status

		statuses, err = phlPost(p.TrackingNumber)
		total = len(statuses)

		p.History = statuses
	}

	return
}

//////////////////////////////////////////////////////////
//
//
//                       Local methods
//
//
/////////////////////////////////////////////////////////
func phlPost(tn string) (statuses []*Status, err error) {
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
	statuses = []*Status{}
	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		// Exclude header
		if i == 0 {
			return
		}

		columns := []string{}
		row.Children().Each(func(i int, col *goquery.Selection) {
			columns = append(columns, strings.Trim(col.Text(), " "))
		})

		// Parse time and subtract offset UTC+8 since PHLPOST
		// is using PH timezone
		t, _ := time.Parse("Jan 02 2006 3:04PM", columns[1])
		t = t.Add(-8 * time.Hour)

		statuses = append(statuses, &Status{t.UTC().Unix(), columns[0], columns[2]})
	})

	return
}
