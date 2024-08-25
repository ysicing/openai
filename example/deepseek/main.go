package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/ysicing/openai/openai"
)

func main() {
	client, err := openai.New(
		openai.WithToken(os.Getenv("DEEPSEEK_API")),
		openai.WithProvider(openai.DEEPSEEK),
		openai.WithModel(openai.DeepseekChat),
	)
	if err != nil {
		panic(err)
	}
	resp, err := client.Completion(context.Background(), "这是博客评论, 请帮我检查是否涉及广告、色情、政治等不宜公开的信息, 并且用'yes'或'no'回答:\n\n 五毛钱一条删除")
	if err == nil {
		log.Printf("prompt:%d,completion:%d,total:%d", resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
		if strings.Contains(resp.Content, "yes") {
			log.Println("rejected")
		} else {
			log.Println("approved")
		}
	}
	resp2, err := client.Completion(context.Background(), "五毛钱一条删除", "这是博客评论, 请帮我检查是否涉及广告、色情、政治等不宜公开的信息, 并且用'yes'或'no'回答:")
	if err == nil {
		log.Printf("prompt:%d,completion:%d,total:%d", resp2.Usage.PromptTokens, resp2.Usage.CompletionTokens, resp2.Usage.TotalTokens)
		if strings.Contains(resp2.Content, "yes") {
			log.Println("rejected")
		} else {
			log.Println("approved")
		}
	}
}
