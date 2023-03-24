//go:build integration

package amazonpolly

import (
	_ "embed"
	"testing"
)

var (
	//go:embed secrets/AWS_ACCESS_KEY_ID.txt
	AWS_ACCESS_KEY_ID string
	//go:embed secrets/AWS_SECRET_ACCESS_KEY.txt
	AWS_SECRET_ACCESS_KEY string
)

func TestAmazonPolly_GenerateAudioFromChaptors(t *testing.T) {
	type args struct {
		fairyTaleChaptors []string
		outputFilename    string
	}
	tests := []struct {
		name    string
		polly   *AmazonPolly
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			polly: &AmazonPolly{
				awsKeyId:           AWS_ACCESS_KEY_ID,
				awsSecretAccessKey: AWS_SECRET_ACCESS_KEY,
			},
			args: args{
				fairyTaleChaptors: []string{"Amy is a princess", "Bob is a prince", "They are in Germany"},
				outputFilename:    `C:\temp\test.mp3`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.polly.GenerateAudioFromChaptors(tt.args.fairyTaleChaptors, tt.args.outputFilename); (err != nil) != tt.wantErr {
				t.Errorf("AmazonPolly.GenerateAudioFromChaptors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
