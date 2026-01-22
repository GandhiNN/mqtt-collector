// Provides automatic payload type detection for MQTT messages
// using hierarchical format analysis
package payload

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"mqtt-catalog/pkg/models"
	"unicode/utf8"
)

// Analyzes byte payload and classifies as JSON (objects/arrays only),
// XML, plain text, or binary data (default)
func DetectType(payload []byte) models.PayloadType {
	if len(payload) == 0 {
		return models.PayloadText
	}

	// JSON validator
	// Only consider objects {} and arrays [] as JSON
	trimmed := bytes.TrimSpace(payload)
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		var js json.RawMessage
		if json.Unmarshal(payload, &js) == nil {
			return models.PayloadJSON
		}
	}

	// XML validator
	if xml.Unmarshal(payload, new(interface{})) == nil {
		return models.PayloadXML
	}

	// UTF-8 text validator
	if utf8.Valid(payload) {
		return models.PayloadText
	}

	return models.PayloadBinary
}
