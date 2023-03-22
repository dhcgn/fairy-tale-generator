package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	
main characters: %s
support characters: %s

the story takes place in %s
the plot of the main characters is %s

The fairy tale should be funny and entertaining for children, it should be about 5-10 pages long.
`, mainCharatersAggregated, supporterCharatersAggregated, location, storyPlot)

	return prompt
}

// generateFairyTaleText generates the fairy tale text with the help of the OpenAI GPT-3 API.
func generateFairyTaleText(apiKey string, mainCharaters []string, supporterCharaters []string, location, storyPlot string) (string, string, error) {
	prompt := generateChatGtpPrompt(mainCharaters, supporterCharaters, location, storyPlot)
	data := map[string]interface{}{
		"model":       "text-davinci-003", // gpt-3.5-turbo or text-davinci-003
		"prompt":      prompt,
		"n":           1,
		"max_tokens":  4097 - len(prompt)/4, // model's maximum context length is 4097 tokens
		"temperature": 0.8,
	}
	// max tokens: https://platform.openai.com/docs/models/gpt-4

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("OpenAI-Organization", "org-K1e47v2URc9rSlglooxQReFe")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var result ResponseOpenAi
	if err = json.Unmarshal(body, &result); err != nil {
		return "", "", err
	}

	if result.Error.Message != "" {
		return "", "", fmt.Errorf("OpenAI API error: %s", result.Error.Message)
	}

	choices := result.Choices
	firstChoice := choices[0]
	text := firstChoice.Text

	return text, prompt, nil
}
