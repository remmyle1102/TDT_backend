package models

// Report represent the report model
type Report struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Date     string `json:"date"`
	AddBy    string `json:"addBy"`
}

type ReportTableData struct {
	Host       string        `json:"host"`
	ReportData []*ReportData `json:"reportData"`
}

type ReportData struct {
	Folder          string                 `json:"folder"`
	File            []string               `json:"file"`
	FileData        []FileData             `json:"fileData"`
	DBTaskData      AuditDBInstance        `json:"dbTaskData"`
	TableSuggestion map[string]interface{} `json:"taskSuggestion"`
	CheckError      map[string]interface{} `json:"checkError"`
}

type FileData struct {
	FileData []map[string]interface{}
}
