package heimat

// Day _
type Day struct {
	Date         string       `json:"date"`
	TrackedTimes []TrackEntry `json:"trackedTimes"`
}

// Month is the same as Day. it just has more Tracked Times
type Month struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	LockStatus struct {
		ID int `json:"id"`
	} `json:"lockStatus"`
	TrackedTimesDate []Day `json:"trackedTimesDate"`
}
