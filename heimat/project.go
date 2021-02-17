package heimat

// Project _
type Project struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

// Task _
type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
