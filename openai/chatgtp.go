package openai

import (
	"bytes"
	"encoding/json"
	"fairy-tale-generator/story"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pterm/pterm"
)

func New(apiKey string, orgID string, model string, story story.FairyTaleOptions) *OpenAI {
	return &OpenAI{
		apiKey: apiKey,
		orgID:  orgID,
		model:  model,
		story:  story,
	}
}

type OpenAI struct {
	apiKey string
	orgID  string
	model  string
	story  story.FairyTaleOptions
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

func (ai *OpenAI) generateFairyTaleTextInternal(r request) (*ChatCompletion, error) {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ai.apiKey))
	req.Header.Set("OpenAI-Organization", ai.orgID)

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
func (ai *OpenAI) GenerateFairyTaleText() ([]string, string, error) {
	prompt := story.GenerateChatGtpPrompt(ai.story)
	conservation := []Message{
		{assistant, prompt},
	}
	data := request{
		Model:    ai.model,
		Messages: conservation,
	}

	chapters := []string{}

	pterm.Info.Printf("Generating %v. chapter with the OpenAI model %v ...\n", 1, ai.model)

	response, err := ai.generateFairyTaleTextInternal(data)

	if err != nil {
		return nil, "", err
	}
	if response.Error.Message != "" {
		return nil, "", fmt.Errorf("OpenAI API error: %s", response.Error.Message)
	}

	chapters = append(chapters, response.Choices[0].Message.Content)
	conservation = append(conservation, response.Choices[0].Message)
	conservation = append(conservation, Message{assistant, "Write next chapter."})

	for i := 0; i < ai.story.ChapterCount-1; i++ {
		pterm.Info.Printf("Generating %v. chapter with the OpenAI model %v ...\n", i+2, ai.model)

		response, err = ai.generateFairyTaleTextInternal(request{
			Model:    ai.model,
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
