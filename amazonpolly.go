package main

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

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
