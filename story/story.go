package story

import (
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
