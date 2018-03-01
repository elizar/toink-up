package parcel

// Constants
const (
	PHLPOST = "phlpost"
	LBC     = "lbc"
)

// Status ...
type Status struct {
	Time     int64  `json:"time"`
	Status   string `json:"status"`
	Location string `json:"location"`
}

// Tracker ...
type Parcel struct {
	Delivered       bool      `json:"delivered"`
	Courier         string    `json:"courier"`
	TrackingNumber  string    `json:"trackingNumber"`
	TrackingHistory []*Status `json:"trackingHistory"`
}
