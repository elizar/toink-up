package parcel

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// NewTracker ...
func NewParcel(courier, tn string) *Tracker {
	return &Parcel{
		Courier:        courier,
		TrackingNumber: tn,
	}
}

func (p *Parcel) Fetch() (total int, err error) {

}

func phlPost(tn string) (parcels []*parcel, err error) {
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
	parcels = []*parcel{}
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

		parcels = append(parcels, &parcel{
			columns[0],
			t.UTC().Unix(),
			columns[2],
			false,
		})
	})

	// Empty
	if len(parcels) == 0 {
		err = errors.New("Package does not exist")
	}

	return
}
