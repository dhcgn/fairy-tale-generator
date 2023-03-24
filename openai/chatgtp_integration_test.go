//go:build integration

package openai

import (
	"os"
	"reflect"
	"strings"
	"testing"

	_ "embed"
)

var (
	//go:embed secrets/OPENAI_API_KEY.txt
	OPENAI_API_KEY string
	//go:embed secrets/OPENAI_ORGANIZATION.txt
	OPENAI_ORGANIZATION string
)

func TestOpenAI_GenerateFairyTaleText(t *testing.T) {
	tests := []struct {
		name    string
		ai      *OpenAI
		want    []string
		want1   string
		wantErr bool
	}{
		{
			name: "test1",
			ai: &OpenAI{
				apiKey: OPENAI_API_KEY,
				orgID:  OPENAI_ORGANIZATION,
				model:  "gpt-3.5-turbo",
				story: FairyTaleOptions{
					MainCharaters:      []string{"Amy"},
					SupporterCharaters: []string{"Bob"},
					Location:           "Germany",
					StoryPlot:          "Amy is a princess",
					ChapterCount:       3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.ai.GenerateFairyTaleText()
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenAI.GenerateFairyTaleText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			os.WriteFile(`C:\temp\test.chaptors.txt`, []byte(strings.Join(got, "\n\n")), os.ModePerm)
			os.WriteFile(`C:\temp\test.prompt.txt`, []byte(got1), os.ModePerm)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenAI.GenerateFairyTaleText() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("OpenAI.GenerateFairyTaleText() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
