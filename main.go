package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "embed"

	"runtime/debug"

	"github.com/joho/godotenv"
	"github.com/pterm/pterm"
)

var (
	// Version is the version of the program.
	Version = "dev"
)

var (
	apiKey       = os.Getenv("OPENAI_API_KEY")
	orgID        = os.Getenv("OPENAI_ORGANIZATION")
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
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
	CharacterSupporterSet = textToLines(characterSupporterSetAsset)
	// LocationSet is the set of locations for the fairy tale.
	LocationSet = textToLines(locationSetAsset)
	// StoryPlotSet is the set of story plots for the fairy tale.
	StoryPlotSet = textToLines(storyPlotsSetAsset)
	// CharacterMainSet is the set of main characters for the fairy tale.
	CharacterMainSet = textToLines(characterMainSetAsset)
)

var (
	flatSetRandomOptions = flag.Bool("random", false, "Use a flat set of random options instead of the default interactive options")
)

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}()

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	fmt.Println("Fairy Tale Generator")
	fmt.Println("Version:", Version, "Commit:", Commit)
	fmt.Println("https://github.com/dhcgn/fairy-tale-generator")
	fmt.Println()
	fmt.Println("This program uses OpenAI's GPT-3 API to generate a fairy tale in German with the help of AWS Polly.")
	fmt.Println()
	pterm.Info.Printf("Loaded %v main charaters, %v support charaters, %v locations, %v plots which results in %v possible stories\n", len(CharacterMainSet), len(CharacterSupporterSet), len(LocationSet), len(StoryPlotSet), len(CharacterMainSet)*len(CharacterSupporterSet)*len(LocationSet)*len(StoryPlotSet))
	fmt.Println()

	flag.Parse()

	if apiKey == "" || awsAccessKey == "" || awsSecretKey == "" || orgID == "" {
		pterm.Error.Println("Please set the environment variables OPENAI_API_KEY, OPENAI_ORGANIZATION, AWS_ACCESS_KEY_ID, and AWS_SECRET_ACCESS_KEY.")
		return
	}

	targetFolder := "results"
	err := os.Mkdir(targetFolder, 0755)
	if err != nil && !os.IsExist(err) {
		pterm.Error.Printf("Error creating results directory: %v\n", err)
		return
	}

	fairyTaleOptions := getFairyTaleOptions(*flatSetRandomOptions)

	generateAndPlay(fairyTaleOptions, targetFolder)
}

// generateAndPlay generates the fairy tale text and audio and plays the audio.
func generateAndPlay(opts fairyTaleOptions, targetFolder string) {
	pterm.Info.Println("Generating fairy tale")

	fairyTaleChaptors, prompt, err := generateFairyTaleText(apiKey, orgID, opts)
	if err != nil {
		pterm.Error.Printf("Error generating fairy tale: %v\n", err)
		return
	}

	ts := createTimestamp()
	f, _ := os.Create(filepath.Join(targetFolder, fmt.Sprintf("%s_fairy_tale.txt", ts)))
	defer f.Close()

	f.WriteString("Prompt:\n")
	f.WriteString(prompt)
	f.WriteString("\n\nPlot:\n")
	f.WriteString(strings.Join(fairyTaleChaptors, "\n\n"))

	pterm.Info.Println("Generated fairy tale saved!")

	pterm.Info.Println("Generating audio from fairy tale")

	outputFilename := filepath.Join(targetFolder, fmt.Sprintf("%s_fairy_tale.mp3", ts))
	generateAudioFromChaptors(fairyTaleChaptors, outputFilename)

	pterm.Info.Printf("German audio saved to: %s\n", outputFilename)

	if runtime.GOOS == "windows" {
		// Open the audio file on windows
		cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", outputFilename)
		err = cmd.Run()
		if err != nil {
			pterm.Error.Printf("Error opening audio file: %v\n", err)
			return
		}
	}
}

func textToLines(text string) []string {
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

func createTimestamp() string {
	now := time.Now().UTC()
	return now.Format("2006-01-02_15-04-05")
}
