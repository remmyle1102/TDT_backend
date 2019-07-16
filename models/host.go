package models

// Host represent the host model
type Host struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IPAdd       string `json:"ipAdd"`
	Port        int    `json:"port"`
	Description string `json:"description"`
	AddBy       string `json:"addBy"`
}
