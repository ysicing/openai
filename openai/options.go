package openai

import (
	"errors"
	"time"

	"github.com/sashabaranov/go-openai"
)

var (
	errorsMissingToken = errors.New("please set OPENAI_API_KEY environment variable")
)

const (
	// DEPRECATED OPENAI is deprecated, use OpenAI instead
	OPENAI = "openai"
	OpenAI = "openai"
	// AZURE is deprecated, use Azure instead
	AZURE = "azure"
	Azure = "azure"
	// DEEPSEEK is deprecated, use DeepSeek instead
	DEEPSEEK = "deepseek"
	DeepSeek = "deepseek"
	ZhiPu    = "zhipu"
)

const (
	defaultMaxTokens   = 2000
	defaultModel       = openai.GPT4oMini
	defaultTemperature = 1.0
	defaultProvider    = OpenAI
	defaultTopP        = 1.0
)

// Option is an interface that specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

// optionFunc is a type of function that can be used to implement the Option interface.
// It takes a pointer to a config struct and modifies it.
type optionFunc func(*config)

// Ensure that optionFunc satisfies the Option interface.
var _ Option = (*optionFunc)(nil)

// The apply method of optionFunc type is implemented here to modify the config struct based on the function passed.
func (o optionFunc) apply(c *config) {
	o(c)
}

// WithToken is a function that returns an Option, which sets the token field of the config struct.
func WithToken(val string) Option {
	return optionFunc(func(c *config) {
		c.token = val
	})
}

// WithOrgID is a function that returns an Option, which sets the orgID field of the config struct.
func WithOrgID(val string) Option {
	return optionFunc(func(c *config) {
		c.orgID = val
	})
}

// WithModel is a function that returns an Option, which sets the model field of the config struct.
func WithModel(val string) Option {
	return optionFunc(func(c *config) {
		c.model = val
	})
}

// WithProxyURL is a function that returns an Option, which sets the proxyURL field of the config struct.
func WithProxyURL(val string) Option {
	return optionFunc(func(c *config) {
		c.proxyURL = val
	})
}

// WithSocksURL is a function that returns an Option, which sets the socksURL field of the config struct.
func WithSocksURL(val string) Option {
	return optionFunc(func(c *config) {
		c.socksURL = val
	})
}

// WithBaseURL returns a new Option that sets the base URL for the client configuration.
// It takes a string value representing the base URL to use for requests.
// It returns an optionFunc that sets the baseURL field of the configuration to the provided
func WithBaseURL(val string) Option {
	return optionFunc(func(c *config) {
		c.baseURL = val
	})
}

// WithTimeout returns a new Option that sets the timeout for the client configuration.
// It takes a time.Duration value representing the timeout duration.
// It returns an optionFunc that sets the timeout field of the configuration to the provided value.
func WithTimeout(val time.Duration) Option {
	return optionFunc(func(c *config) {
		c.timeout = val
	})
}

// WithMaxTokens returns a new Option that sets the max tokens for the client configuration.
// The maximum number of tokens to generate in the chat completion.
// The total length of input tokens and generated tokens is limited by the model's context length.
func WithMaxTokens(val int) Option {
	if val <= 0 {
		val = defaultMaxTokens
	}
	return optionFunc(func(c *config) {
		c.maxTokens = val
	})
}

// WithTemperature returns a new Option that sets the temperature for the client configuration.
// What sampling temperature to use, between 0 and 2.
// Higher values like 0.8 will make the output more random,
// while lower values like 0.2 will make it more focused and deterministic.
func WithTemperature(val float32) Option {
	if val <= 0 {
		val = defaultTemperature
	}
	return optionFunc(func(c *config) {
		c.temperature = val
	})
}

// WithProvider sets the `provider` variable based on the value of the `val` parameter.
// If `val` is not set to `OpenAI` or `Azure`, it will be set to the default value `defaultProvider`.
// This function returns an `Option` object.
func WithProvider(val string) Option {
	// Check if `val` is set to `OpenAI` or `Azure`. If not, set it to the default value.
	switch val {
	case OpenAI, Azure, DeepSeek, ZhiPu:
	default:
		val = defaultProvider
	}

	// Return an `optionFunc` object with `c.provider` set to `val`.
	return optionFunc(func(c *config) {
		c.provider = val
	})
}

// WithSkipVerify returns a new Option that sets the skipVerify for the client configuration.
func WithSkipVerify(val bool) Option {
	return optionFunc(func(c *config) {
		c.skipVerify = val
	})
}

// WithHeaders returns a new Option that sets the headers for the http client configuration.
func WithHeaders(headers []string) Option {
	return optionFunc(func(c *config) {
		c.headers = headers
	})
}

// WithApiVersion returns a new Option that sets the apiVersion for OpenAI Model.
func WithApiVersion(apiVersion string) Option {
	return optionFunc(func(c *config) {
		c.apiVersion = apiVersion
	})
}

// WithTopP returns a new Option that sets the topP for the client configuration.
func WithTopP(val float32) Option {
	return optionFunc(func(c *config) {
		c.topP = val
	})
}

// WithPresencePenalty returns a new Option that sets the presencePenalty for the client configuration.
func WithPresencePenalty(val float32) Option {
	return optionFunc(func(c *config) {
		c.presencePenalty = val
	})
}

// WithFrequencyPenalty returns a new Option that sets the frequencyPenalty for the client configuration.
func WithFrequencyPenalty(val float32) Option {
	return optionFunc(func(c *config) {
		c.frequencyPenalty = val
	})
}

// config is a struct that stores configuration options for the instrumentation.
type config struct {
	baseURL     string
	token       string
	orgID       string
	model       string
	proxyURL    string
	socksURL    string
	timeout     time.Duration
	maxTokens   int
	temperature float32

	topP             float32
	presencePenalty  float32
	frequencyPenalty float32

	provider   string
	skipVerify bool
	headers    []string
	apiVersion string
}

// valid checks whether a config object is valid, returning an error if it is not.
func (cfg *config) valid() error {
	// Check that the token is not empty.
	if cfg.token == "" {
		return errorsMissingToken
	}

	if cfg.provider == DEEPSEEK || cfg.provider == DeepSeek {
		cfg.model = DeepseekChat
		return nil
	}

	if cfg.provider == ZhiPu {
		if len(cfg.model) == 0 {
			cfg.model = ZhiPuGlmFree
		}
		return nil
	}

	if (cfg.provider == OpenAI || cfg.provider == Azure) && len(cfg.model) == 0 {
		cfg.model = defaultModel
	}
	// If all checks pass, return nil (no error).
	return nil
}

// newConfig creates a new config object with default values, and applies the given options.
func newConfig(opts ...Option) *config {
	// Create a new config object with default values.
	c := &config{
		maxTokens:   defaultMaxTokens,
		temperature: defaultTemperature,
		provider:    defaultProvider,
		topP:        defaultTopP,
	}

	// Apply each of the given options to the config object.
	for _, opt := range opts {
		opt.apply(c)
	}

	// Return the resulting config object.
	return c
}
