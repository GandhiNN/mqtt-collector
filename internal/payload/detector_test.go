// Tests the payload type detection logic:
// - Empty payloads
// - Various JSON formats (object, array, string, number)
// - XML documents
// - Plain text
// - Binary data
// - Invalid UTF-8
// - Benchmark test for performance testing

package payload

import (
	"mqtt-collector/pkg/models"
	"testing"
)

func TestDetectType(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		expected models.PayloadType
	}{
		{
			name:     "empty payload",
			payload:  []byte{},
			expected: models.PayloadText,
		},
		{
			name:     "valid JSON object",
			payload:  []byte(`{"temperature": 22.5, "humidity": 60}`),
			expected: models.PayloadJSON,
		},
		{
			name:     "valid JSON array",
			payload:  []byte(`[1, 2, 3, 4, 5]`),
			expected: models.PayloadJSON,
		},
		{
			name:     "valid JSON string as text",
			payload:  []byte(`{"hello world"}`),
			expected: models.PayloadText,
		},
		{
			name:     "valid JSON number as text",
			payload:  []byte(`{42}`),
			expected: models.PayloadText,
		},
		{
			name:     "valid XML",
			payload:  []byte(`<?xml version="1.0"?><root><item>value</item></root>`),
			expected: models.PayloadXML,
		},
		{
			name:     "simple XML",
			payload:  []byte(`<message>Hello</message>`),
			expected: models.PayloadXML,
		},
		{
			name:     "plain text",
			payload:  []byte(`This is plain text without any structure`),
			expected: models.PayloadText,
		},
		{
			name:     "number as text",
			payload:  []byte(`12345`),
			expected: models.PayloadText,
		},
		{
			name:     "binary data",
			payload:  []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}, // JPEG header
			expected: models.PayloadBinary,
		},
		{
			name:     "invalid UTF-8",
			payload:  []byte{0x80, 0x81, 0x82, 0x83},
			expected: models.PayloadBinary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectType(tt.payload)
			if result != tt.expected {
				t.Errorf(
					"DetectType() = %v, want %v for payload: %q",
					result,
					tt.expected,
					tt.payload,
				)
			}
		})
	}
}

func BenchmarkDetectType(b *testing.B) {
	payloads := [][]byte{
		[]byte(`{"temperature": 22.5}`),
		[]byte(`<root><item>value</item></root>`),
		[]byte(`plain text message`),
		{0xFF, 0xDB, 0xFF, 0xE0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectType(payloads[i%len(payloads)])
	}
}
