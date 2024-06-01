package storage

import "time"

// easyjson:json
type Event struct {
	ID          string        `json:"id,omitempty"`
	Title       string        `json:"title"`
	StartTime   time.Time     `json:"startTime"`
	EndTime     time.Time     `json:"endTime"`
	Description string        `json:"description,omitempty"`
	OwnerID     string        `json:"ownerId"`
	Notify      time.Duration `json:"notify,omitempty"`
}
