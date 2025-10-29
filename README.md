# OpenAI Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/ysicing/openai)](https://goreportcard.com/report/github.com/ysicing/openai)
[![Test Coverage](https://img.shields.io/badge/coverage-71.2%25-brightgreen)](https://github.com/ysicing/openai)
[![Go Version](https://img.shields.io/badge/go-1.24-blue)](https://golang.org/dl/)

A **flexible and secure** Go SDK for OpenAI and OpenAI-compatible APIs including **local models** like **Ollama**.

## ✨ Features

- 🔌 **Multi-Provider Support**: OpenAI, Azure OpenAI, DeepSeek, ZhiPu, **Ollama**, and any OpenAI-compatible API
- 🛡️ **Security First**: Comprehensive validation, TLS warnings, safe error handling
- 🧪 **Well-Tested**: 71.2% test coverage with comprehensive unit tests
- 🚀 **Developer-Friendly**: Clean functional options pattern, extensive documentation
- ⚡ **Performant**: Optimized HTTP client, connection pooling, context-aware
- 🔧 **Flexible**: Support for custom BaseURL, proxies, timeouts, and all OpenAI parameters

## 📦 Installation

```bash
go get github.com/ysicing/openai
```

## 🚀 Quick Start

### OpenAI
```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/ysicing/openai/openai"
)

func main() {
    client, err := openai.New(
        openai.WithToken(os.Getenv("OPENAI_API_KEY")),
        openai.WithModel(openai.GPT4oMini),
    )
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Completion(
        context.Background(),
        "You are a helpful assistant.",
        "Explain quantum computing in simple terms",
    )
    if err != nil {
        log.Fatal(err)
    }

    log.Println(resp.Content)
}
```

### Ollama (Local AI)
```go
client, err := openai.New(
    openai.WithToken("ollama"),
    openai.WithProvider(openai.Ollama),  // 自动配置 localhost:11434
    openai.WithModel("llama3.1:latest"),
    openai.WithTimeout(300),  // 本地模型需要更长超时
)
```

### DeepSeek (OpenAI-Compatible)
```go
client, err := openai.New(
    openai.WithToken(os.Getenv("DEEPSEEK_API_KEY")),
    openai.WithBaseURL("https://api.deepseek.com/v1"),
    openai.WithModel(openai.DeepseekChat),
)
```

## 🔧 Configuration Options

| Option | Description | Example |
|--------|-------------|---------|
| `WithToken` | API authentication token | `"sk-..."` |
| `WithModel` | Model name | `"gpt-4o-mini"` |
| `WithProvider` | Service provider | `openai.Ollama` |
| `WithBaseURL` | Custom API endpoint | `"http://localhost:11434/v1"` |
| `WithMaxTokens` | Maximum output tokens | `2000` |
| `WithTemperature` | Response creativity (0-2) | `0.7` |
| `WithTopP` | Nucleus sampling | `0.9` |
| `WithTimeout` | Request timeout | `60 * time.Second` |
| `WithProxyURL` | HTTP proxy | `"http://proxy:8080"` |
| `WithSocksURL` | SOCKS5 proxy | `"socks5://proxy:1080"` |
| `WithSkipVerify` | Skip TLS verification | `true` ⚠️ |

## 📚 Supported Providers

### Built-in Providers (Special Configuration)

| Provider | Official Support | Default URL | Notes |
|----------|-----------------|-------------|-------|
| OpenAI | ✅ | `https://api.openai.com/v1` | Default provider |
| Azure OpenAI | ✅ | Azure endpoint | Enterprise features |
| **Ollama** | ✅ | `http://localhost:11434/v1` | **Local models** |

✅ = Built-in provider with special configuration

### OpenAI-Compatible Providers (Via WithBaseURL)

All other providers use the **default OpenAI-compatible mode**. Just specify the endpoint:

```go
// DeepSeek
openai.WithBaseURL("https://api.deepseek.com/v1")

// ZhiPu (智谱)
openai.WithBaseURL("https://open.bigmodel.cn/api/paas/v4/")

// LM Studio, LocalAI, vLLM, etc.
openai.WithBaseURL("http://localhost:1234/v1")
```

**Compatible with any OpenAI-compatible service:**
- DeepSeek, ZhiPu, Moonshot, Kimi, Qwen, etc. (Chinese providers)
- LM Studio, LocalAI, vLLM (Self-hosted)
- Any other service with OpenAI-compatible API

## 💡 Design Philosophy

We use a **simplified approach**:
- ✅ Only OpenAI, Azure, and Ollama have **special configurations**
- ✅ All other providers use **default OpenAI-compatible mode**
- ✅ Just use `WithBaseURL` to specify the endpoint

This design makes the code simpler and supports more services!

## 🎯 Advanced Usage

### Multi-turn Conversations
```go
messages := []openai.ChatCompletionMessage{
    {Role: openai.ChatMessageRoleUser, Content: "Hello!"},
}

resp, err := client.CreateChatCompletionWithMessage(context.Background(), messages)
if err != nil {
    log.Fatal(err)
}

messages = append(messages, resp.Choices[0].Message)
messages = append(messages, openai.ChatCompletionMessage{
    Role: openai.ChatMessageRoleUser,
    Content: "What is your name?",
})

resp2, err := client.CreateChatCompletionWithMessage(context.Background(), messages)
```

### Image Understanding (GPT-4V)
```go
resp, err := client.ImageCompletion(
    context.Background(),
    "https://example.com/image.jpg",
    "You are a helpful assistant.",
    "Describe this image in detail",
)
```

### Custom Headers
```go
client, err := openai.New(
    openai.WithToken("api-key"),
    openai.WithHeaders([]string{
        "X-Custom-Header=custom-value",
        "Authorization=Bearer token",
    }),
)
```

## 🛡️ Security Features

- ✅ **Response Validation**: Prevents panics on malformed API responses
- ✅ **TLS Warnings**: Clear documentation about security implications of `WithSkipVerify`
- ✅ **Safe Error Handling**: No credential leakage in error messages
- ✅ **Input Validation**: Robust header parsing with `SplitN`
- ✅ **Proxy Security**: Support for both HTTP and SOCKS5 proxies

## 🧪 Testing

Run all tests:
```bash
go test ./... -v
```

With coverage report:
```bash
go test ./... -cover
```

Test results:
```
PASS
coverage: 71.2% of statements
ok      github.com/ysicing/openai/openai    0.496s
```

## 📖 Examples

See the `example/` directory for complete working examples:
- `example/deepseek/` - DeepSeek API usage
- `example/zhipu/` - ZhiPu API usage
- `example/ollama/` - Local Ollama usage

```bash
# Run DeepSeek example
cd example/deepseek
DEEPSEEK_API_KEY=your_key go run main.go

# Run Ollama example (requires Ollama running)
cd example/ollama
ollama run llama3.1:latest
go run main.go
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

- [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) - OpenAI Go library
- [Ollama](https://ollama.ai) - Local AI models
- All contributors and users of this SDK
