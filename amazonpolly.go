package main

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/hyacinthus/mp3join"
	"github.com/pterm/pterm"
)

func generateAudioFromChaptors(fairyTaleChaptors []string, outputFilename string) error {
	joiner := mp3join.New()

	for i, fairyTaleChaptor := range fairyTaleChaptors {
		pterm.Info.Println("Converting text to speech Chapter", i+1, "of", len(fairyTaleChaptors))
		tempAudioFile := fmt.Sprintf("%v_temp.%v.mp3", outputFilename, i)
		err := generateAudioAndSaveToDisk(fairyTaleChaptor, tempAudioFile)
		if err != nil {
			pterm.Error.Printf("Error converting text to speech: %v\n", err)
			return err
		}

		f_temp, err := os.Open(tempAudioFile)
		if err != nil {
			pterm.Error.Printf("Error opening audio file: %v\n", err)
			return err
		}
		err = joiner.Append(f_temp)
		if err != nil {
			pterm.Error.Printf("Error appending audio file: %v\n", err)
			return err
		}
		f_temp.Close()
		err = os.Remove(tempAudioFile)
		if err != nil {
			pterm.Error.Printf("Error removing audio file: %v\n", err)
			return err
		}
	}

	f, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err != nil {
		pterm.Error.Printf("Error creating audio file: %v\n", err)
		return err
	}
	_, err = io.Copy(f, joiner.Reader())
	if err != nil {
		pterm.Error.Printf("Error copying audio file from joiner: %v\n", err)
		return err
	}
	return nil
}

// generateAudioAndSaveToDisk generates an audio file from the given text and saves it to disk
func generateAudioAndSaveToDisk(text, outputFilename string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := polly.New(sess, aws.NewConfig().WithRegion("eu-central-1"))

	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		SampleRate:   aws.String("22050"),
		Text:         aws.String(text),
		TextType:     aws.String("text"),
		VoiceId:      aws.String("Daniel"), // German voice
		Engine:       aws.String("neural"),
	}

	output, err := svc.SynthesizeSpeech(input)
	if err != nil {
		return err
	}

	defer output.AudioStream.Close()

	buf, err := io.ReadAll(output.AudioStream)
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFilename, buf, 0644)
	if err != nil {
		return err
	}

	return nil
}
