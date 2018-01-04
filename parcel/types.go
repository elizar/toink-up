package parcel

// Status ...
type Status struct {
	Status   string `json:"location"`
	Time     int64  `json:"time"`
	Location string `json:"location"`
}

// Tracker ...
type Parcel struct {
	Courier        string    `json:"courier"`
	TrackingNumber string    `json:"trackingNumber"`
	Delivered      bool      `json:"delivered"`
	History        []*Status `json:"history"`
}
