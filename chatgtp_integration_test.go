//go:build integration

package main

import (
	_ "embed"
	"testing"
)

var (
	//go:embed secrets/OPENAI_API_KEY.txt
	OPENAI_API_KEY string
	//go:embed secrets/OPENAI_ORGANIZATION.txt
	OPENAI_ORGANIZATION string
)

func Test_generateFairyTaleText(t *testing.T) {
	type args struct {
		apiKey             string
		orgID              string
		mainCharaters      []string
		supporterCharaters []string
		location           string
		storyPlot          string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want1   string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				apiKey:             OPENAI_API_KEY,
				orgID:              OPENAI_ORGANIZATION,
				mainCharaters:      []string{"Amy"},
				supporterCharaters: []string{"Bob"},
				location:           "Germany",
				storyPlot:          "Amy is a princess",
			},
			want:  nil,
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := generateFairyTaleText(tt.args.apiKey, tt.args.orgID, tt.args.mainCharaters, tt.args.supporterCharaters, tt.args.location, tt.args.storyPlot)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateFairyTaleText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
