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
)

// DefaultModel is the default OpenAI model to use if one is not provided.
var DefaultModel = openai.GPT4oMini

// modelMaps maps model names to their corresponding model ID strings.
var modelMaps = map[string]string{
	"gpt-4":          openai.GPT4,
	"gpt-4o-mini":    openai.GPT4oMini,
	"gpt-4-turbo":    openai.GPT4Turbo,
	"gpt-3.5-turbo":  openai.GPT3Dot5Turbo,
	"deepseek-chat":  DeepseekChat,
	"deepseek-coder": DeepseekCoder,
}

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
		model:       modelMaps[cfg.model],
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
	case AZURE:
		// Set the OpenAI client to use the default configuration with Azure-specific options, if the provider is Azure.
		defaultAzureConfig := openai.DefaultAzureConfig(cfg.token, cfg.baseURL)
		defaultAzureConfig.AzureModelMapperFunc = func(model string) string {
			return cfg.modelName
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
	case DEEPSEEK:
		{
			c.HTTPClient = httpClient
			if cfg.apiVersion != "" {
				c.APIVersion = cfg.apiVersion
			}
			if cfg.baseURL == "" {
				c.BaseURL = "https://api.deepseek.com"
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
	content ...string,
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
				Role:    openai.ChatMessageRoleUser,
				Content: content[0],
			},
		},
	}

	if len(content) > 1 {
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: content[1],
		})
	}

	return c.client.CreateChatCompletion(ctx, req)
}

// CreateCompletion is an API call to create a completion.
// This is the main endpoint of the API. It returns new text, as well as, if requested,
// the probabilities over each alternative token at each position.
//
// If using a fine-tuned model, simply provide the model's ID in the CompletionRequest object,
// and the server will use the model's parameters to generate the completion.
func (c *Client) CreateCompletion(
	ctx context.Context,
	content string,
) (resp openai.CompletionResponse, err error) {
	req := openai.CompletionRequest{
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		FrequencyPenalty: c.frequencyPenalty,
		PresencePenalty:  c.presencePenalty,
		Prompt:           content,
	}

	return c.client.CreateCompletion(ctx, req)
}

// Completion is a method on the Client struct that takes a context.Context and a string argument
// and returns a string and an error.
func (c *Client) Completion(
	ctx context.Context,
	content ...string,
) (*Response, error) {
	resp := &Response{}
	switch c.model {
	case openai.GPT3Dot5Turbo,
		openai.GPT3Dot5Turbo0301,
		openai.GPT3Dot5Turbo0613,
		openai.GPT3Dot5Turbo1106,
		openai.GPT3Dot5Turbo0125,
		openai.GPT3Dot5Turbo16K,
		openai.GPT3Dot5Turbo16K0613,
		openai.GPT4,
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
		openai.GPT432K0613,
		DeepseekChat,
		DeepseekCoder:
		r, err := c.CreateChatCompletion(ctx, content...)
		if err != nil {
			return nil, err
		}
		resp.Content = r.Choices[0].Message.Content
		resp.Usage = r.Usage
	default:
		r, err := c.CreateCompletion(ctx, content[0])
		if err != nil {
			return nil, err
		}
		resp.Content = r.Choices[0].Text
		resp.Usage = r.Usage
	}
	return resp, nil
}

// CreateImageChatCompletion is an API call to create a completion for a chat message.
func (c *Client) CreateImageChatCompletion(
	ctx context.Context,
	image string,
	content ...string,
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
						Text: content[0],
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

	if len(content) > 1 {
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: content[1],
		})
	}

	return c.client.CreateChatCompletion(ctx, req)
}

// ImageCompletion is a method on the Client struct that takes a context.Context and a string argument
func (c *Client) ImageCompletion(
	ctx context.Context,
	image string,
	content ...string,
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
		r, err := c.CreateImageChatCompletion(ctx, image, content...)
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
