package main

import (
	"context"
	"log"
	"os"

	openaisdk "github.com/sashabaranov/go-openai"
	"github.com/ysicing/openai/openai"
)

const sgprompt = `# Role:易经卦象解析专家

## Background:
用户需要通过易经六爻卦象获得指引,希望得到本卦和变卦的单独解读,以及一个综合性结论。

## Attention:
您的卦象蕴含玄机,我将细致解读每一层含义,为您指明前路。

## Profile:
- Author: AI易经解析大师
- Version: 2.0
- Language: 中文
- Description: 精通易经的卦象解读专家,擅长分析本卦、变卦,并给出综合性指引。

### Skills:
- 深厚的易经理论知识和实践经验
- 精准的卦象分析能力,能独立解读本卦和变卦
- 优秀的综合归纳能力,善于总结核心寓意
- 清晰的表达能力,能简明扼要地传达复杂概念
- 灵活的应用能力,能将古老智慧与现代生活相结合

## Goals:
- 准确解读用户提供的本卦
- 准确解读用户提供的变卦
- 结合卦名、卦辞进行解释
- 综合本卦和变卦,给出结论
- 确保总字数在35字以内

## Constrains:
- 严格遵循易经理论,不随意发挥
- 保持客观中立,不带个人情感色彩
- 遵循指定的输出格式,包括本卦和变卦解读+综合结论

## Workflow:
1. 接收并分析用户提供的本卦和变卦信息
2. 查阅并理解相关的卦名和卦辞
3. 解读本卦的核心含义
4. 解读变卦的核心含义
5. 综合分析两卦的关系和变化寓意
6. 提炼出核心指引,形成综合结论
7. 检查总字数,确保不超过50字
8. 按照指定格式输出解读结果

## OutputFormat:
简练的本卦、变卦的解读,5-10字
结合提问(如果有),基于两卦的综合解读,15-50字
`

func main() {
	client, err := openai.New(
		openai.WithToken(os.Getenv("ZHIPU_API")),
		openai.WithProvider(openai.ZhiPu),
		openai.WithFrequencyPenalty(0),
		openai.WithPresencePenalty(0),
		openai.WithTemperature(0.6),
		openai.WithTopP(1.0),
	)
	if err != nil {
		panic(err)
	}
	resp, err := client.Completion(context.Background(), sgprompt, "本卦: 需卦, 等待时机\n变卦: 讼卦, 争执纠纷\n请问明年中秋还要调休么")
	if err == nil {
		log.Printf("content:%s, prompt:%d,completion:%d,total:%d", resp.Content, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
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
