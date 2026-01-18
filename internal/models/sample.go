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
