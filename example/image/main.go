package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ysicing/openai/openai"
)

func main() {
	// 示例 1: 使用 OpenAI GPT-4V 进行图像理解
	fmt.Println("=== OpenAI GPT-4V ===")
	openAIClient, err := openai.New(
		openai.WithToken("your-openai-api-key"),
		openai.WithModel("gpt-4o"), // 或其他支持图像的模型
	)
	if err != nil {
		log.Printf("OpenAI client creation failed: %v", err)
	} else {
		resp, err := openAIClient.ImageCompletion(
			context.Background(),
			"https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg",
			"You are a helpful assistant.",
			"Describe this image in detail",
		)
		if err != nil {
			log.Printf("OpenAI error: %v", err)
		} else {
			fmt.Println(resp.Content)
		}
	}

	// 示例 2: 使用 ZhiPu GLM-4V 进行图像理解（如果 ZhiPu 支持图像）
	fmt.Println("\n=== ZhiPu GLM-4V ===")
	zhiPuClient, err := openai.New(
		openai.WithToken("your-zhipu-api-key"),
		openai.WithBaseURL("https://open.bigmodel.cn/api/paas/v4/"),
		openai.WithModel(openai.ZhiPuGlmFree), // 或其他支持图像的 ZhiPu 模型
	)
	if err != nil {
		log.Printf("ZhiPu client creation failed: %v", err)
	} else {
		resp, err := zhiPuClient.ImageCompletion(
			context.Background(),
			"https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg",
			"You are a helpful assistant.",
			"Describe this image in detail",
		)
		if err != nil {
			log.Printf("ZhiPu error: %v", err)
		} else {
			fmt.Println(resp.Content)
		}
	}

	// 示例 3: 使用本地 Ollama + LLaVA 进行图像理解
	fmt.Println("\n=== Ollama + LLaVA (本地) ===")
	ollamaClient, err := openai.New(
		openai.WithToken("ollama"),
		openai.WithProvider(openai.Ollama),
		openai.WithModel("llava:latest"), // LLaVA 模型
	)
	if err != nil {
		log.Printf("Ollama client creation failed: %v", err)
	} else {
		resp, err := ollamaClient.ImageCompletion(
			context.Background(),
			"https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg",
			"You are a helpful assistant.",
			"Describe this image in detail",
		)
		if err != nil {
			log.Printf("Ollama error: %v", err)
		} else {
			fmt.Println(resp.Content)
		}
	}

	// 示例 4: 使用 LM Studio + Vision 模型
	fmt.Println("\n=== LM Studio + Vision Model (本地) ===")
	lmStudioClient, err := openai.New(
		openai.WithToken("lm-studio"),
		openai.WithBaseURL("http://localhost:1234/v1"),
		openai.WithModel("lmstudio-community/Llava-Vision-Model"), // 任何支持视觉的模型
	)
	if err != nil {
		log.Printf("LM Studio client creation failed: %v", err)
	} else {
		resp, err := lmStudioClient.ImageCompletion(
			context.Background(),
			"https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg",
			"You are a helpful assistant.",
			"Describe this image in detail",
		)
		if err != nil {
			log.Printf("LM Studio error: %v", err)
		} else {
			fmt.Println(resp.Content)
		}
	}
}
