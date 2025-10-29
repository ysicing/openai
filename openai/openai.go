package openai

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/net/proxy"
)

const (
	// Common model names for convenience
	// These work with any OpenAI-compatible provider
	DeepseekChat = "deepseek-chat"
	ZhiPuGlmFree = "glm-4-flash"
)

// DefaultModel is the default OpenAI model to use if one is not provided.
var DefaultModel = openai.GPT4oMini

// Client is a struct that represents an OpenAI client.
type Client struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32

	// An alternative to sampling with temperature, called nucleus sampling,
	// where the model considers the results of the tokens with top_p probability mass.
	// So 0.1 means only the tokens comprising the top 10% probability mass are considered.
	topP float32
	// Number between -2.0 and 2.0.
	// Positive values penalize new tokens based on whether they appear in the text so far,
	// increasing the model's likelihood to talk about new topics.
	presencePenalty float32
	// Number between -2.0 and 2.0.
	// Positive values penalize new tokens based on their existing frequency in the text so far,
	// decreasing the model's likelihood to repeat the same line verbatim.
	frequencyPenalty float32
}

type Response struct {
	Content string
	Usage   openai.Usage
}

// New creates a new OpenAI API client with the given options.
func New(opts ...Option) (*Client, error) {
	// Create a new config object with the given options.
	cfg := newConfig(opts...)

	// Validate the config object, returning an error if it is invalid.
	if err := cfg.valid(); err != nil {
		return nil, err
	}

	// Create a new client instance with the necessary fields.
	engine := &Client{
		model:       cfg.model,
		maxTokens:   cfg.maxTokens,
		temperature: cfg.temperature,
	}

	// Create a new OpenAI config object with the given API token and other optional fields.
	c := openai.DefaultConfig(cfg.token)
	if cfg.orgID != "" {
		c.OrgID = cfg.orgID
	}
	if cfg.baseURL != "" {
		c.BaseURL = cfg.baseURL
	}

	// Create a new HTTP transport.
	tr := &http.Transport{}
	if cfg.skipVerify {
		// WARNING: Disabling TLS certificate verification exposes the client to
		// man-in-the-middle (MITM) attacks. This should ONLY be used in development
		// or testing environments. Never use this in production with sensitive data.
		// Consider using proper certificates instead.
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Create a new HTTP client with the specified timeout and proxy, if any.
	httpClient := &http.Client{
		Timeout: cfg.timeout,
	}

	if cfg.proxyURL != "" {
		proxyURL, _ := url.Parse(cfg.proxyURL)
		tr.Proxy = http.ProxyURL(proxyURL)
	} else if cfg.socksURL != "" {
		dialer, err := proxy.SOCKS5("tcp", cfg.socksURL, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("proxy connection failed: verify SOCKS5 proxy address and network connectivity")
		}
		tr.DialContext = dialer.(proxy.ContextDialer).DialContext
	}

	// Set the HTTP client to use the default header transport with the specified headers.
	httpClient.Transport = &DefaultHeaderTransport{
		Origin: tr,
		Header: NewHeaders(cfg.headers),
	}

	switch cfg.provider {
	case Azure:
		// Azure OpenAI has special configuration requirements
		defaultAzureConfig := openai.DefaultAzureConfig(cfg.token, cfg.baseURL)
		defaultAzureConfig.AzureModelMapperFunc = func(model string) string {
			return cfg.model
		}
		if cfg.apiVersion != "" {
			defaultAzureConfig.APIVersion = cfg.apiVersion
		}
		defaultAzureConfig.HTTPClient = httpClient
		engine.client = openai.NewClientWithConfig(defaultAzureConfig)

	default:
		// Default mode: OpenAI-compatible API
		// This works for OpenAI, Ollama, DeepSeek, ZhiPu, LM Studio, LocalAI, vLLM, etc.
		c.HTTPClient = httpClient
		if cfg.apiVersion != "" {
			c.APIVersion = cfg.apiVersion
		}
		engine.client = openai.NewClientWithConfig(c)
	}
	// Return the resulting client engine.
	return engine, nil
}

// buildChatCompletionRequest creates a standardized chat completion request
// with common configuration parameters.
func (c *Client) buildChatCompletionRequest(
	messages []openai.ChatCompletionMessage,
) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		FrequencyPenalty: c.frequencyPenalty,
		PresencePenalty:  c.presencePenalty,
		Messages:         messages,
	}
}

// CreateChatCompletion is an API call to create a completion for a chat message.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	prompt,
	content string,
) (resp openai.ChatCompletionResponse, err error) {
	if len(prompt) == 0 {
		prompt = "You are a helpful assistant."
	}
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		},
	}

	req := c.buildChatCompletionRequest(messages)
	return c.client.CreateChatCompletion(ctx, req)
}

// CreateChatCompletionWithMessage is an API call to create a completion for a chat message.
func (c *Client) CreateChatCompletionWithMessage(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
) (resp openai.ChatCompletionResponse, err error) {
	req := c.buildChatCompletionRequest(messages)
	return c.client.CreateChatCompletion(ctx, req)
}

// Completion is a method on the Client struct that takes a context.Context and a string argument
// and returns a Response and an error.
func (c *Client) Completion(
	ctx context.Context,
	prompt, content string,
) (*Response, error) {
	resp := &Response{}
	r, err := c.CreateChatCompletion(ctx, prompt, content)
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	// Validate response to prevent panics on empty choices
	if len(r.Choices) == 0 {
		return nil, errors.New("empty response from API: no choices returned")
	}

	resp.Content = r.Choices[0].Message.Content
	resp.Usage = r.Usage
	return resp, nil
}

// CreateImageChatCompletion is an API call to create a completion for a chat message with image input.
func (c *Client) CreateImageChatCompletion(
	ctx context.Context,
	image, prompt, content string,
) (resp openai.ChatCompletionResponse, err error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeText,
					Text: content,
				},
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL: image,
					},
				},
			},
		},
	}
	if len(prompt) > 0 {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		})
	}

	req := c.buildChatCompletionRequest(messages)
	return c.client.CreateChatCompletion(ctx, req)
}

// ImageCompletion is a method on the Client struct for image understanding.
// It sends an image and text prompt to the model and returns the response.
// This method works with any model that supports image input (e.g., GPT-4V, Claude 3, etc.)
// The underlying OpenAI client will validate if the model supports image input.
func (c *Client) ImageCompletion(
	ctx context.Context,
	image, prompt, content string,
) (*Response, error) {
	r, err := c.CreateImageChatCompletion(ctx, image, prompt, content)
	if err != nil {
		return nil, fmt.Errorf("image chat completion failed: %w", err)
	}

	// Validate response to prevent panics on empty choices
	if len(r.Choices) == 0 {
		return nil, errors.New("empty response from API: no choices returned")
	}

	return &Response{
		Content: r.Choices[0].Message.Content,
		Usage:   r.Usage,
	}, nil
}
