package main

import (
	_ "embed"
	"os"
	"testing"
)

var (
	//go:embed secrets/AWS_ACCESS_KEY_ID.txt
	AWS_ACCESS_KEY_ID string
	//go:embed secrets/AWS_SECRET_ACCESS_KEY.txt
	AWS_SECRET_ACCESS_KEY string
)

func Test_generateAudioFromChaptors(t *testing.T) {

	os.Setenv("AWS_ACCESS_KEY_ID", AWS_ACCESS_KEY_ID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", AWS_SECRET_ACCESS_KEY)

	type args struct {
		fairyTaleChaptors []string
		outputFilename    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				fairyTaleChaptors: []string{"Amy is a princess", "Bob is a prince", "They are in Germany"},
				outputFilename:    `C:\temp\test.mp3`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := generateAudioFromChaptors(tt.args.fairyTaleChaptors, tt.args.outputFilename); (err != nil) != tt.wantErr {
				t.Errorf("generateAudioFromChaptors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
