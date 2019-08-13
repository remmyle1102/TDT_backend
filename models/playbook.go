package models

// Playbook represent the playbook model
type Playbook struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	AddBy       string     `json:"addBy"`
	Description NullString `json:"description"`
	Location    string     `json:"location"`
}

// PlaybookContent represent the playbook content
type PlaybookContent struct {
	PlaybookContent string `json:"playbookContent"`
}

type PlaybookTemplate struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	SubTasks    []map[string]interface{} `json:"subTasks"`
}
