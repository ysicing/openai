package main

import (
	"context"
	"log"
	"os"
	"strings"

	openaisdk "github.com/sashabaranov/go-openai"
	"github.com/ysicing/openai/openai"
)

func main() {
	client, err := openai.New(
		openai.WithToken(os.Getenv("ZHIPU_API")),
		openai.WithProvider(openai.ZhiPu),
		openai.WithModel(openai.ZhiPuGlmFree),
	)
	if err != nil {
		panic(err)
	}
	resp, err := client.Completion(context.Background(), "你是资深审核家", "这是博客评论, 请帮我检查是否涉及广告、色情、政治等不宜公开的信息, 并且用'yes'或'no'回答:\n\n 五毛钱一条删除")
	if err == nil {
		log.Printf("content:%s, prompt:%d,completion:%d,total:%d", resp.Content, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
		if strings.Contains(resp.Content, "yes") {
			log.Printf("rejected")
		} else {
			log.Printf("approved")
		}
	}
	messages := []openaisdk.ChatCompletionMessage{
		{
			Role:    openaisdk.ChatMessageRoleUser,
			Content: "What's the highest mountain in the world?",
		},
	}
	resp2, err := client.CreateChatCompletionWithMessage(context.Background(), messages)
	if err != nil {
		log.Printf("error:%v", err)
		return
	}
	log.Printf("content:%s, prompt:%d,completion:%d,total:%d", resp2.Choices[0].Message.Content, resp2.Usage.PromptTokens, resp2.Usage.CompletionTokens, resp2.Usage.TotalTokens)
	messages = append(messages, resp2.Choices[0].Message)
	// spew.Dump(messages)
	messages = append(messages, openaisdk.ChatCompletionMessage{
		Role:    openaisdk.ChatMessageRoleUser,
		Content: "What is the second?",
	})
	resp3, err := client.CreateChatCompletionWithMessage(context.Background(), messages)
	if err != nil {
		log.Printf("error:%v", err)
		return
	}
	log.Printf("content:%s, prompt:%d,completion:%d,total:%d", resp3.Choices[0].Message.Content, resp3.Usage.PromptTokens, resp3.Usage.CompletionTokens, resp3.Usage.TotalTokens)
	// messages = append(messages, resp3.Choices[0].Message)
	// spew.Dump(messages)
}
