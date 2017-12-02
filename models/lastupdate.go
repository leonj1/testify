package models

type LastUpdate struct {
	Date   MyTime `json:"date,omitempty"`
	By     string `json:"by,omitempty"`
	Status string `json:"status,omitempty"`
}
