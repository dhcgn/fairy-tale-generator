package story

import (
	"fmt"
	"math/rand"
	"strings"

	_ "embed"
)

// Embeded text files
var (
	//go:embed asssets/charactersupporterset.txt
	characterSupporterSetAsset string
	//go:embed asssets/locationset.txt
	locationSetAsset string
	//go:embed asssets/storyplotsset.txt
	storyPlotsSetAsset string
	//go:embed asssets/charactermainset.txt
	characterMainSetAsset string
)

// Options for the fairy tale
var (
	// CharacterSupporterSet is the set of supporter characters for the fairy tale.
	CharacterSupporterSet = textToLinesInRandomOrder(characterSupporterSetAsset)
	// LocationSet is the set of locations for the fairy tale.
	LocationSet = textToLinesInRandomOrder(locationSetAsset)
	// StoryPlotSet is the set of story plots for the fairy tale.
	StoryPlotSet = textToLinesInRandomOrder(storyPlotsSetAsset)
	// CharacterMainSet is the set of main characters for the fairy tale.
	CharacterMainSet = textToLinesInRandomOrder(characterMainSetAsset)
)

type FairyTaleOptions struct {
	MainCharaters      []string
	SupporterCharaters []string
	Location           string
	StoryPlot          string
	ChapterCount       int
}

func textToLinesInRandomOrder(text string) []string {
	lines := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line != "" {
			lines = append(lines, strings.Trim(line, ""))
		}
	}

	rand.Shuffle(len(lines), func(i, j int) {
		lines[i], lines[j] = lines[j], lines[i]
	})

	return lines
}

// generateChatGtpPrompt generates the prompt for the chat gtp
func GenerateChatGtpPrompt(story FairyTaleOptions) string {
	mainCharatersAggregated := aggregateSlice(story.MainCharaters)
	supporterCharatersAggregated := aggregateSlice(story.SupporterCharaters)

	prompt := fmt.Sprintf(`Write a long children fairy tale in German.
	
Main characters: %s
Support characters: %s

The story takes place in %s
The plot of the main characters is: %s

The fairy tale should be funny, entertaining for children.
Write it in %v chapters and start only with the first chapter.
`, mainCharatersAggregated, supporterCharatersAggregated, story.Location, story.StoryPlot, story.ChapterCount)

	return prompt
}

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
