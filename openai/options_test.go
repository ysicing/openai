package openai

import (
	"testing"
	"time"
)

func TestWithToken(t *testing.T) {
	opt := WithToken("test-token")

	c := &config{}
	opt.apply(c)

	if c.token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", c.token)
	}
}

func TestWithOrgID(t *testing.T) {
	opt := WithOrgID("test-org")

	c := &config{}
	opt.apply(c)

	if c.orgID != "test-org" {
		t.Errorf("Expected orgID 'test-org', got '%s'", c.orgID)
	}
}

func TestWithModel(t *testing.T) {
	opt := WithModel("test-model")

	c := &config{}
	opt.apply(c)

	if c.model != "test-model" {
		t.Errorf("Expected model 'test-model', got '%s'", c.model)
	}
}

func TestWithTemperature(t *testing.T) {
	// Test with valid value
	opt := WithTemperature(0.5)
	c := &config{}
	opt.apply(c)

	if c.temperature != 0.5 {
		t.Errorf("Expected temperature 0.5, got %f", c.temperature)
	}

	// Test with zero value (should use default)
	opt = WithTemperature(0)
	c = &config{}
	opt.apply(c)

	if c.temperature != defaultTemperature {
		t.Errorf("Expected default temperature %f, got %f", defaultTemperature, c.temperature)
	}
}

func TestWithTimeout(t *testing.T) {
	duration := 30 * time.Second
	opt := WithTimeout(duration)

	c := &config{}
	opt.apply(c)

	if c.timeout != duration {
		t.Errorf("Expected timeout %v, got %v", duration, c.timeout)
	}
}

func TestWithProxyURL(t *testing.T) {
	opt := WithProxyURL("http://proxy.example.com")

	c := &config{}
	opt.apply(c)

	if c.proxyURL != "http://proxy.example.com" {
		t.Errorf("Expected proxyURL 'http://proxy.example.com', got '%s'", c.proxyURL)
	}
}

func TestWithSocksURL(t *testing.T) {
	opt := WithSocksURL("socks5://proxy.example.com")

	c := &config{}
	opt.apply(c)

	if c.socksURL != "socks5://proxy.example.com" {
		t.Errorf("Expected socksURL 'socks5://proxy.example.com', got '%s'", c.socksURL)
	}
}

func TestWithBaseURL(t *testing.T) {
	opt := WithBaseURL("https://custom.api.com")

	c := &config{}
	opt.apply(c)

	if c.baseURL != "https://custom.api.com" {
		t.Errorf("Expected baseURL 'https://custom.api.com', got '%s'", c.baseURL)
	}
}

func TestWithProvider(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"OpenAI", OpenAI, OpenAI},
		{"Azure", Azure, Azure},
		{"Ollama (uses default)", "ollama", defaultProvider},     // Uses default mode
		{"DeepSeek (uses default)", "deepseek", defaultProvider}, // Uses default mode
		{"ZhiPu (uses default)", "zhipu", defaultProvider},       // Uses default mode
		{"Unknown provider", "unknown", defaultProvider},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithProvider(tt.input)
			c := &config{}
			opt.apply(c)

			if c.provider != tt.expected {
				t.Errorf("Expected provider '%s', got '%s'", tt.expected, c.provider)
			}
		})
	}
}

func TestWithSkipVerify(t *testing.T) {
	opt := WithSkipVerify(true)

	c := &config{}
	opt.apply(c)

	if !c.skipVerify {
		t.Error("Expected skipVerify to be true")
	}
}

func TestWithHeaders(t *testing.T) {
	headers := []string{"Authorization=Bearer token", "Content-Type=application/json"}
	opt := WithHeaders(headers)

	c := &config{}
	opt.apply(c)

	if len(c.headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(c.headers))
	}
}

func TestWithApiVersion(t *testing.T) {
	opt := WithApiVersion("v1")

	c := &config{}
	opt.apply(c)

	if c.apiVersion != "v1" {
		t.Errorf("Expected apiVersion 'v1', got '%s'", c.apiVersion)
	}
}

func TestWithTopP(t *testing.T) {
	opt := WithTopP(0.9)

	c := &config{}
	opt.apply(c)

	if c.topP != 0.9 {
		t.Errorf("Expected topP 0.9, got %f", c.topP)
	}
}

func TestWithPresencePenalty(t *testing.T) {
	opt := WithPresencePenalty(0.5)

	c := &config{}
	opt.apply(c)

	if c.presencePenalty != 0.5 {
		t.Errorf("Expected presencePenalty 0.5, got %f", c.presencePenalty)
	}
}

func TestWithFrequencyPenalty(t *testing.T) {
	opt := WithFrequencyPenalty(0.3)

	c := &config{}
	opt.apply(c)

	if c.frequencyPenalty != 0.3 {
		t.Errorf("Expected frequencyPenalty 0.3, got %f", c.frequencyPenalty)
	}
}

func TestConfig_Valid(t *testing.T) {
	// Test with valid token
	c := &config{
		token:    "test-token",
		provider: OpenAI,
	}

	err := c.valid()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test with missing token
	c = &config{
		token: "",
	}

	err = c.valid()
	if err == nil {
		t.Error("Expected error for missing token, got nil")
	}

	// Test OpenAI provider sets default model
	c = &config{
		token:    "test-token",
		provider: OpenAI,
	}

	err = c.valid()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if c.model != defaultModel {
		t.Errorf("Expected default model to be set, got '%s'", c.model)
	}

	// Test custom provider (uses default mode)
	c = &config{
		token:    "test-token",
		provider: "custom-provider",
		model:    "custom-model",
	}

	err = c.valid()
	if err != nil {
		t.Errorf("Expected no error for custom provider, got: %v", err)
	}
}

func TestNewConfig(t *testing.T) {
	// Test with no options
	c := newConfig()

	if c.temperature != defaultTemperature {
		t.Errorf("Expected default temperature %f, got %f", defaultTemperature, c.temperature)
	}

	if c.provider != defaultProvider {
		t.Errorf("Expected default provider '%s', got '%s'", defaultProvider, c.provider)
	}

	if c.topP != defaultTopP {
		t.Errorf("Expected default topP %f, got %f", defaultTopP, c.topP)
	}

	// Test with options
	c = newConfig(
		WithToken("test-token"),
		WithModel("test-model"),
		WithTemperature(0.7),
		WithProvider("custom-provider"),
	)

	if c.token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", c.token)
	}

	if c.model != "test-model" {
		t.Errorf("Expected model 'test-model', got '%s'", c.model)
	}

	if c.temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", c.temperature)
	}

	// Unknown providers default to OpenAI-compatible mode
	if c.provider != defaultProvider {
		t.Errorf("Expected provider '%s', got '%s'", defaultProvider, c.provider)
	}
}

func TestOptionInterface(t *testing.T) {
	// Test that optionFunc implements Option interface
	var _ Option = (*optionFunc)(nil)

	// Test that apply method works
	opt := optionFunc(func(c *config) {
		c.token = "test-token"
	})

	c := &config{}
	opt.apply(c)

	if c.token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", c.token)
	}
}
