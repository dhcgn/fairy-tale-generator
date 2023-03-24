package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pterm/pterm"
)

func aggregateSlice(input []string) string {
	var output string
	for i, v := range input {
		if i == len(input)-1 {
			output += v
		} else {
			output += fmt.Sprintf("%s, ", v)
		}
	}
	return output
}

// generateChatGtpPrompt generates the prompt for the chat gtp
func generateChatGtpPrompt(mainCharaters []string, supporterCharaters []string, location, storyPlot string) string {
	mainCharatersAggregated := aggregateSlice(mainCharaters)
	supporterCharatersAggregated := aggregateSlice(supporterCharaters)

	prompt := fmt.Sprintf(`Create a children fairy tale in German with the following setup.
	
Main characters: %s
Support characters: %s

The story takes place in %s
The plot of the main characters is: %s

The fairy tale should be funny, entertaining for children and in german.
Write it in %v chapters and start only with the first chapter.
`, mainCharatersAggregated, supporterCharatersAggregated, location, storyPlot, ChapterCount)

	return prompt
}

type request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	// MaxTokens int       `json:"max_tokens"`
}

// OpenAI API roles
const (
	system    = "system"
	assistant = "assistant"
	user      = "user"
)

func generateFairyTaleTextInternal(apiKey string, r request) (*ChatCompletion, error) {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("OpenAI-Organization", orgID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ChatCompletion
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// generateFairyTaleText generates the fairy tale text with the help of the OpenAI GPT-3 API.
func generateFairyTaleText(apiKey string, orgID string, opts fairyTaleOptions) ([]string, string, error) {
	prompt := generateChatGtpPrompt(opts.mainCharaters, opts.supporterCharaters, opts.location, opts.storyPlot)
	conservation := []Message{
		{assistant, prompt},
	}
	data := request{
		Model:    "gpt-3.5-turbo",
		Messages: conservation,
	}
	// max tokens: https://platform.openai.com/docs/models/gpt-4

	chapters := []string{}

	pterm.Info.Println("Generating 1. chapter ...")

	response, err := generateFairyTaleTextInternal(apiKey, data)

	if err != nil {
		return nil, "", err
	}
	if response.Error.Message != "" {
		return nil, "", fmt.Errorf("OpenAI API error: %s", response.Error.Message)
	}

	chapters = append(chapters, response.Choices[0].Message.Content)
	conservation = append(conservation, response.Choices[0].Message)
	conservation = append(conservation, Message{assistant, "Write next chapter."})

	for i := 0; i < ChapterCount-1; i++ {
		pterm.Info.Printf("Generating %v. chapter ...\n", i+2)

		response, err = generateFairyTaleTextInternal(apiKey, request{
			Model:    "gpt-3.5-turbo",
			Messages: conservation,
		})
		if err != nil {
			return nil, "", err
		}
		if response.Error.Message != "" {
			return nil, "", fmt.Errorf("OpenAI API error: %s", response.Error.Message)
		}
		chapters = append(chapters, response.Choices[0].Message.Content)
		conservation = append(conservation, response.Choices[0].Message)
		conservation = append(conservation, Message{assistant, "Write next chapter."})
	}

	pterm.Info.Printf("Story generated with %v words!\n", len(strings.Fields(strings.Join(chapters, "\n"))))

	return chapters, prompt, nil
}

type ChatCompletion struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Choices []Choice  `json:"choices"`
	Usage   UsageInfo `json:"usage"`
	Error   struct {
		Message string      `json:"message"`
		Type    string      `json:"type"`
		Param   interface{} `json:"param"`
		Code    interface{} `json:"code"`
	} `json:"error"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
