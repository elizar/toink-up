package parcel

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
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
	var statuses []*Status

	switch p.Courier {
	case PHLPOST:
		statuses, err = phlPost(p.TrackingNumber)
	case LBC:
		statuses, err = lbc(p.TrackingNumber)
	default:
		err = errors.New("Courier does not exists")
		return
	}

	total = len(statuses)

	if total == 0 {
		err = errors.New("package not found")
		return
	}

	p.TrackingHistory = statuses
	rx := regexp.MustCompile(`(?i)item delivered`)
	p.Delivered = rx.MatchString(p.TrackingHistory[0].Status)

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

	statuses = []*Status{}

	// Loop through each row
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

		statuses = append(statuses, &Status{t.Unix(), columns[0], columns[2]})
	})

	return
}

func lbc(tn string) (statuses []*Status, err error) {
	const endpoint = "http://www.lbcexpress.com/track/?"

	// Params
	p := make(url.Values)
	p.Add("tracking_no", tn) // "1122039085"

	resp, err := http.PostForm(endpoint, p)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}

	statuses = []*Status{}

	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		columns := []string{}
		row.Children().Each(func(i int, col *goquery.Selection) {
			columns = append(columns, strings.Trim(col.Text(), " "))
		})

		t, _ := time.Parse("January 2, 2006", columns[0])

		statuses = append(statuses, &Status{t.Unix(), columns[1], columns[2]})
	})

	return
}
