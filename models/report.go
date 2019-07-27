package models

// Report represent the report model
type Report struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ProblemCount int    `json:"problemCount"`
	Location     string `json:"location"`
	Date         string `json:"date"`
	AddBy        string `json:"addBy"`
}

type ReportData struct {
	Folder     string          `json:"folder"`
	File       []string        `json:"file"`
	FileData   []FileData      `json:"fileData"`
	DBTaskData AuditDBInstance `json:"dbTaskData"`
}

type FileData struct {
	FileData []map[string]interface{}
}
