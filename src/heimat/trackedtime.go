package heimat

// TrackedTime _
type TrackedTime struct {
	ID       int      `json:"id"`
	Start    string   `json:"start"`
	End      string   `json:"end"`
	Date     string   `json:"date"`
	Note     string   `json:"note"`
	Project  Project  `json:"project"`
	Task     Task     `json:"task"`
	Employee Employee `json:"employee"`
}
