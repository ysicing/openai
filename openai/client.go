package openai

import (
	"net/http"
	"strings"
)

// DefaultHeaderTransport is an http.RoundTripper that adds the given headers to
type DefaultHeaderTransport struct {
	Origin http.RoundTripper
	Header http.Header
}

// RoundTrip implements the http.RoundTripper interface.
func (t *DefaultHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, values := range t.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return t.Origin.RoundTrip(req)
}

// NewHeaders creates a new http.Header from the given slice of headers.
// Headers should be in the format "Key=Value". Malformed headers are skipped.
func NewHeaders(headers []string) http.Header {
	h := make(http.Header)
	for _, header := range headers {
		// Split header into key and value with = as delimiter, limit to 2 parts
		// Using SplitN prevents panic if header contains multiple '=' characters
		vals := strings.SplitN(header, "=", 2)
		if len(vals) != 2 {
			continue
		}

		// Trim whitespace from key and value
		key := strings.TrimSpace(vals[0])
		value := strings.TrimSpace(vals[1])

		// Skip empty keys
		if key == "" {
			continue
		}

		h.Add(key, value)
	}
	return h
}
