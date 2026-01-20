package models

import "time"

type PayloadType string

const (
	PayloadJSON   PayloadType = "json"
	PayloadXML    PayloadType = "xml"
	PayloadText   PayloadType = "text"
	PayloadBinary PayloadType = "binary"
)

type Sample struct {
	BrokerID    string      `json:"broker_id"`
	Topic       string      `json:"topic"`
	PayloadType PayloadType `json:"payload_type"`
	Payload     []byte      `json:"payload"`
	Timestamp   time.Time   `json:"timestamp"`
}

// Topic is the database model
type Topic struct {
	ID            int64       `json:"id"             db:"id"`
	BrokerID      string      `json:"broker_id"      db:"broker_id"`
	Topic         string      `json:"topic"          db:"topic"`
	PayloadType   PayloadType `json:"payload_type"   db:"payload_type"`
	SamplePayload []byte      `json:"sample_payload" db:"sample_payload"`
	LastSeen      time.Time   `json:"last_seen"      db:"last_seen"`
	CreatedAt     time.Time   `json:"created_at"     db:"created_at"`
}

type TopicListResponse struct {
	Topics []Topic `json:"topics"`
	Total  int     `json:"total"`
}
