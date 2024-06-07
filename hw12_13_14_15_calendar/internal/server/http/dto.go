package internalhttp

import "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage"

// easyjson:json
type dto struct {
	Ok      bool            `json:"ok"`
	Details string          `json:"details,omitempty"`
	Count   int             `json:"count"`
	Events  []storage.Event `json:"items,omitempty"`
}
