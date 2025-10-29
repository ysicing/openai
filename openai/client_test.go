package openai

import (
	"net/http"
	"strings"
	"testing"
)

func TestNewHeaders(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string][]string
	}{
		{
			name:  "Valid headers",
			input: []string{"Authorization=Bearer token", "Content-Type=application/json"},
			expected: map[string][]string{
				"Authorization": {"Bearer token"},
				"Content-Type":  {"application/json"},
			},
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: map[string][]string{},
		},
		{
			name:  "Header with spaces",
			input: []string{" Authorization = Bearer token ", "Content-Type=application/json"},
			expected: map[string][]string{
				"Authorization": {"Bearer token"},
				"Content-Type":  {"application/json"},
			},
		},
		{
			name:     "Invalid header (no equals)",
			input:    []string{"InvalidHeader"},
			expected: map[string][]string{},
		},
		{
			name:     "Invalid header (empty key)",
			input:    []string{"=value"},
			expected: map[string][]string{},
		},
		{
			name:  "Multiple equals signs",
			input: []string{"Cookie=value=with=equals"},
			expected: map[string][]string{
				"Cookie": {"value=with=equals"},
			},
		},
		{
			name:  "Empty value",
			input: []string{"Header="},
			expected: map[string][]string{
				"Header": {""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewHeaders(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d headers, got %d", len(tt.expected), len(result))
			}

			for key, expectedValues := range tt.expected {
				actualValues := result[key]
				if len(actualValues) != len(expectedValues) {
					t.Errorf("Expected %d values for key '%s', got %d",
						len(expectedValues), key, len(actualValues))
				}

				for i, expectedValue := range expectedValues {
					if i >= len(actualValues) {
						t.Errorf("Missing value at index %d for key '%s'", i, key)
						continue
					}

					if actualValues[i] != expectedValue {
						t.Errorf("Expected value '%s' for key '%s', got '%s'",
							expectedValue, key, actualValues[i])
					}
				}
			}
		})
	}
}

func TestDefaultHeaderTransport_RoundTrip(t *testing.T) {
	origin := &mockRoundTripper{}
	headers := make(http.Header)
	headers.Add("X-Custom-Header", "test-value")

	transport := &DefaultHeaderTransport{
		Origin: origin,
		Header: headers,
	}

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	_, err = transport.RoundTrip(req)
	if err != nil {
		t.Errorf("RoundTrip failed: %v", err)
	}

	// Verify that the custom header was added
	if origin.request.Header.Get("X-Custom-Header") == "" {
		t.Error("Expected X-Custom-Header to be added to request")
	}

	expectedValue := origin.request.Header.Get("X-Custom-Header")
	if expectedValue != "test-value" {
		t.Errorf("Expected header value 'test-value', got '%s'", expectedValue)
	}
}

type mockRoundTripper struct {
	request *http.Request
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.request = req
	return &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}, nil
}

func TestNewHeaders_InputValidation(t *testing.T) {
	// Test that NewHeaders doesn't panic on malformed input
	// This verifies our fix for the string splitting vulnerability

	malformedInputs := [][]string{
		{},
		{"="},
		{"=value"},
		{"key==value"},
		{"key=value1=value2"},
		{"key"},
		{"key="},
		{"   "},
		{strings.Repeat("=", 100)},
	}

	for i, input := range malformedInputs {
		t.Run("Malformed input", func(t *testing.T) {
			// This should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic on input %d: %v", i, r)
				}
			}()

			result := NewHeaders(input)
			_ = result // If we get here, no panic occurred
		})
	}
}
