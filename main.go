package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "embed"

	"github.com/pterm/pterm"
)

var (
	apiKey       = os.Getenv("OPENAI_API_KEY")
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
)

// Embeded text files
var (
	//go:embed asssets\charactersupporterset.txt
	characterSupporterSetAsset string
	//go:embed asssets\locationset.txt
	locationSetAsset string
	//go:embed asssets\storyplotsset.txt
	storyPlotsSetAsset string
	//go:embed asssets\charactermainset.txt
	characterMainSetAsset string
)

// Options for the fairy tale
var (
	// CharacterSupporterSet is the set of supporter characters for the fairy tale.
	CharacterSupporterSet = textToLines(characterSupporterSetAsset)
	// LocationSet is the set of locations for the fairy tale.
	LocationSet = textToLines(locationSetAsset)
	// StoryPlotsSet is the set of story plots for the fairy tale.
	StoryPlotsSet = textToLines(storyPlotsSetAsset)
	// CharacterMainSet is the set of main characters for the fairy tale.
	CharacterMainSet = textToLines(characterMainSetAsset)
)

func main() {
	fmt.Println("Fairy Tale Generator")
	fmt.Println("https://github.com/dhcgn/fairy-tale-generator")
	fmt.Println()
	fmt.Println("This program uses OpenAI's GPT-3 API to generate a fairy tale in German with the help of AWS Polly.")
	fmt.Println()
	pterm.Info.Printf("Loaded %v main charaters, %v support charaters, %v locations, %v plots which results in %v possible stories\n", len(CharacterMainSet), len(CharacterSupporterSet), len(LocationSet), len(StoryPlotsSet), len(CharacterMainSet)*len(CharacterSupporterSet)*len(LocationSet)*len(StoryPlotsSet))
	fmt.Println()

	if apiKey == "" || awsAccessKey == "" || awsSecretKey == "" {
		pterm.Error.Println("Please set the environment variables OPENAI_API_KEY, AWS_ACCESS_KEY_ID, and AWS_SECRET_ACCESS_KEY.")
		return
	}

	selectedMainCharaters, selectedSupporterCharaters, location, storyPlot := getFairyTaleOptions()

	pterm.Info.Println("Generating fairy tale...")

	ts := createTimestamp()
	generateAndPlay(selectedMainCharaters, selectedSupporterCharaters, location, storyPlot, ts)
}

// generateAndPlay generates the fairy tale text and audio and plays the audio.
func generateAndPlay(mainCharaters, supporterCharaters []string, location, storyPlot, ts string) {
	fairyTale, prompt, err := generateFairyTaleText(apiKey, mainCharaters, supporterCharaters, location, storyPlot)
	if err != nil {
		pterm.Error.Printf("Error generating fairy tale: %v\n", err)
		return
	}

	f, _ := os.Create(fmt.Sprintf("%s_fairy_tale.txt", ts))
	defer f.Close()

	f.WriteString("Prompt:\n")
	f.WriteString(prompt)
	f.WriteString("\n\nPlot:\n")
	f.WriteString(fairyTale)

	fmt.Println("Generated fairy tale saved!")

	outputFilename := fmt.Sprintf("%s_fairy_tale.mp3", ts)
	err = generateAudioAndSaveToDisk(fairyTale, outputFilename)
	if err != nil {
		pterm.Error.Printf("Error converting text to speech: %v\n", err)
		return
	}

	fmt.Printf("German audio saved to: %s\n", outputFilename)

	// Open the audio file on windows
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", outputFilename)
	err = cmd.Run()
	if err != nil {
		pterm.Error.Printf("Error opening audio file: %v\n", err)
		return
	}
}

func textToLines(text string) []string {
	lines := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func createTimestamp() string {
	now := time.Now().UTC()
	return now.Format("2006-01-02_15-04-05")
}
