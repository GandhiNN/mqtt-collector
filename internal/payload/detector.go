package payload

import (
	"encoding/json"
	"encoding/xml"
	"mqtt-collector/internal/models"
	"unicode/utf8"
)

func DetectType(payload []byte) models.PayloadType {
	if len(payload) == 0 {
		return models.PayloadText
	}

	// try JSON
	var js json.RawMessage
	if json.Unmarshal(payload, &js) == nil {
		return models.PayloadJSON
	}

	// try XML
	if xml.Unmarshal(payload, new(interface{})) == nil {
		return models.PayloadXML
	}

	// check if valid UTF-8 text
	if utf8.Valid(payload) {
		return models.PayloadText
	}

	return models.PayloadBinary
}
