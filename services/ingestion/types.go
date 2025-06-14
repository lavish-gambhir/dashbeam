package ingestion

import (
	"time"

	"github.com/lavish-gambhir/dashbeam/shared/streaming"
)

type BatchEventsRequest struct {
	Events []streaming.Event `json:"events"`
}

type SingleEventRequest struct {
	Event streaming.Event `json:"event"`
}

type EventResponse struct {
	Status    string    `json:"status"`
	EventIDs  []string  `json:"event_ids,omitempty"`
	Processed int       `json:"processed"`
	Timestamp time.Time `json:"timestamp"`
}
