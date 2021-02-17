package heimat

// Day _
type Day struct {
	Date         string       `json:"date"`
	TrackedTimes []TrackEntry `json:"trackedTimes"`
}
