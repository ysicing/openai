package openai

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/net/proxy"
)

const (
	DeepseekChat  = "deepseek-chat"
	DeepseekCoder = "deepseek-coder"
	ZhiPuGlmFree  = "glm-4-flash"
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
			return nil, fmt.Errorf("can't connect to the proxy: %s", err)
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
		// Set the OpenAI client to use the default configuration with Azure-specific options, if the provider is Azure.
		defaultAzureConfig := openai.DefaultAzureConfig(cfg.token, cfg.baseURL)
		defaultAzureConfig.AzureModelMapperFunc = func(model string) string {
			return cfg.model
		}
		// Set the API version to the one with the specified options.
		if cfg.apiVersion != "" {
			defaultAzureConfig.APIVersion = cfg.apiVersion
		}
		// Set the HTTP client to the one with the specified options.
		defaultAzureConfig.HTTPClient = httpClient
		engine.client = openai.NewClientWithConfig(
			defaultAzureConfig,
		)
	case DeepSeek:
		{
			c.HTTPClient = httpClient
			if cfg.baseURL == "" {
				c.BaseURL = "https://api.deepseek.com"
			}
			engine.client = openai.NewClientWithConfig(c)
		}
	case ZhiPu:
		{
			c.HTTPClient = httpClient
			if cfg.baseURL == "" {
				c.BaseURL = "https://open.bigmodel.cn/api/paas/v4/"
			}
			engine.client = openai.NewClientWithConfig(c)
		}
	default:
		// Otherwise, set the OpenAI client to use the HTTP client with the specified options.
		c.HTTPClient = httpClient
		if cfg.apiVersion != "" {
			c.APIVersion = cfg.apiVersion
		}
		engine.client = openai.NewClientWithConfig(c)
	}
	// Return the resulting client engine.
	return engine, nil
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
	req := openai.ChatCompletionRequest{
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		FrequencyPenalty: c.frequencyPenalty,
		PresencePenalty:  c.presencePenalty,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}
	return c.client.CreateChatCompletion(ctx, req)
}

// CreateChatCompletionWithMessage is an API call to create a completion for a chat message.
func (c *Client) CreateChatCompletionWithMessage(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
) (resp openai.ChatCompletionResponse, err error) {
	req := openai.ChatCompletionRequest{
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		FrequencyPenalty: c.frequencyPenalty,
		PresencePenalty:  c.presencePenalty,
		Messages:         messages,
	}
	return c.client.CreateChatCompletion(ctx, req)
}

// Completion is a method on the Client struct that takes a context.Context and a string argument
// and returns a string and an error.
func (c *Client) Completion(
	ctx context.Context,
	prompt, content string,
) (*Response, error) {
	resp := &Response{}
	r, err := c.CreateChatCompletion(ctx, prompt, content)
	if err != nil {
		return nil, err
	}
	resp.Content = r.Choices[0].Message.Content
	resp.Usage = r.Usage
	return resp, nil
}

// CreateImageChatCompletion is an API call to create a completion for a chat message.
func (c *Client) CreateImageChatCompletion(
	ctx context.Context,
	image, prompt, content string,
) (resp openai.ChatCompletionResponse, err error) {
	req := openai.ChatCompletionRequest{
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		FrequencyPenalty: c.frequencyPenalty,
		PresencePenalty:  c.presencePenalty,
		Messages: []openai.ChatCompletionMessage{
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
		},
	}
	if len(prompt) > 0 {
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		})
	}
	return c.client.CreateChatCompletion(ctx, req)
}

// ImageCompletion is a method on the Client struct that takes a context.Context and a string argument
func (c *Client) ImageCompletion(
	ctx context.Context,
	image, prompt, content string,
) (*Response, error) {
	resp := &Response{}
	switch c.model {
	case openai.GPT4,
		openai.GPT4o,
		openai.GPT4o20240513,
		openai.GPT4o20240806,
		openai.GPT4oMini,
		openai.GPT4oMini20240718,
		openai.GPT4TurboPreview,
		openai.GPT4VisionPreview,
		openai.GPT4Turbo1106,
		openai.GPT4Turbo0125,
		openai.GPT4Turbo,
		openai.GPT4Turbo20240409,
		openai.GPT40314,
		openai.GPT40613,
		openai.GPT432K,
		openai.GPT432K0314,
		openai.GPT432K0613:
		r, err := c.CreateImageChatCompletion(ctx, image, prompt, content)
		if err != nil {
			return nil, err
		}
		resp.Content = r.Choices[0].Message.Content
		resp.Usage = r.Usage
	default:
		return nil, fmt.Errorf("model %s does not support image completions", c.model)
	}
	return resp, nil
}
