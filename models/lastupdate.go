package models

import "time"

type LastUpdate struct {
	Date   time.Time `json:"date,omitempty"`
	By     string    `json:"by,omitempty"`
	Status string    `json:"status,omitempty"`
}
