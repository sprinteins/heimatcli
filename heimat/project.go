package heimat

// Project _
type Project struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

// TaskByName _
func (p Project) TaskByName(name string) *Task {
	for _, task := range p.Tasks {
		if task.Name == name {
			return &task
		}
	}
	return nil
}

// Task _
type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
