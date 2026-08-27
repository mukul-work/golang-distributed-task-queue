package models

import "encoding/json"

type JobRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
