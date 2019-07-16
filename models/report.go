package models

import "time"

// Report represent the report model
type Report struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	ProblemCount int       `json:"ipAddr"`
	Link         string    `json:"link"`
	Date         time.Time `json:"date"`
	AddBy        int       `json:"addBy"`
}
