package parcel

// Constants
const (
	PHLPOST = "phlpost"
)

// Status ...
type Status struct {
	Time     int64  `json:"time"`
	Status   string `json:"status"`
	Location string `json:"location"`
}

// Tracker ...
type Parcel struct {
	Courier        string    `json:"courier"`
	TrackingNumber string    `json:"trackingNumber"`
	Delivered      bool      `json:"delivered"`
	History        []*Status `json:"history"`
}
