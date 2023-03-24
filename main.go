package main

import (
	"fairy-tale-generator/amazonpolly"
	"fairy-tale-generator/openai"
	"fairy-tale-generator/story"
	"flag"
	"fmt"
	"log"
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
	// Version of this program, set during compile time.
	Version = "dev"
)

const (
	// ChapterCount is the number of chaptors in the fairy tale.
	ChapterCount = 3
	// OpenAI API model
	// Most capable GPT-3.5 model and optimized for chat at 1/10th the cost of text-davinci-003. Will be updated with our latest model iteration.
	// Max tokens: 4,096 tokens
	model = "gpt-3.5-turbo"
)

// Secrets for the APIs
var (
	apiKey             = os.Getenv("OPENAI_API_KEY")
	orgID              = os.Getenv("OPENAI_ORGANIZATION")
	awsKeyId           = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
)

var (
	flatSetRandomOptions = flag.Bool("random", false, "Use a set of random options for the strory instead of user interactive input")
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
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	fmt.Println("Fairy Tale Generator")
	fmt.Println("Version:", Version, "Commit:", Commit)
	fmt.Println("https://github.com/dhcgn/fairy-tale-generator")
	fmt.Println()
	fmt.Println("This program uses OpenAI's GPT-3 API to generate a fairy tale in German with the help of AWS Polly.")
	fmt.Println()
	pterm.Info.Printf("Loaded %v main charaters, %v support charaters, %v locations, %v plots which results in %v possible stories\n", len(story.CharacterMainSet), len(story.CharacterSupporterSet), len(story.LocationSet), len(story.StoryPlotSet), len(story.CharacterMainSet)*len(story.CharacterSupporterSet)*len(story.LocationSet)*len(story.StoryPlotSet))
	fmt.Println()

	flag.Parse()

	if apiKey == "" || awsKeyId == "" || awsSecretAccessKey == "" || orgID == "" {
		pterm.Error.Println("Please set the environment variables OPENAI_API_KEY, OPENAI_ORGANIZATION, AWS_ACCESS_KEY_ID, and AWS_SECRET_ACCESS_KEY.")
		return
	}

	targetFolder := "results"
	err := os.Mkdir(targetFolder, 0755)
	if err != nil && !os.IsExist(err) {
		pterm.Error.Printf("Error creating results directory: %v\n", err)
		return
	}

	fairyTaleOptions := story.GetFairyTaleOptions(*flatSetRandomOptions, ChapterCount)

	openai := openai.New(apiKey, orgID, model, fairyTaleOptions)
	amazonpolly := amazonpolly.New(awsKeyId, awsSecretAccessKey)

	generateAndPlay(targetFolder, openai, amazonpolly)
}

// generateAndPlay generates the fairy tale text and audio and plays the audio.
func generateAndPlay(targetFolder string, openai *openai.OpenAI, amazonpolly *amazonpolly.AmazonPolly) {
	pterm.Info.Println("Generating fairy tale")

	fairyTaleChaptors, prompt, err := openai.GenerateFairyTaleText()
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
	amazonpolly.GenerateAudioFromChaptors(fairyTaleChaptors, outputFilename)

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

func createTimestamp() string {
	now := time.Now()
	return now.Format("2006-01-02_15-04-05")
}
