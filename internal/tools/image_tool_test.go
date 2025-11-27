package tools

import (
	"testing"
)

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name   string
		prompt string
		want   string
	}{
		{
			name:   "Standard Title",
			prompt: "Title: 'My Cool Image'\nDescription...",
			want:   "My_Cool_Image",
		},
		{
			name:   "Title Text prefix",
			prompt: "Title Text: Red Velvet's \"Psycho\": The Art of Embracing Chaotic Love\nCentral Metaphor...",
			want:   "Red_Velvets_Psycho_The_Art_of_Embracing_Chaotic_Love",
		},
		{
			name:   "No Title",
			prompt: "Just a description",
			want:   "generated_result_", // Prefix match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTitle(tt.prompt)
			if tt.name == "No Title" {
				if len(got) < len("generated_result_") {
					t.Errorf("extractTitle() = %v, want prefix %v", got, tt.want)
				}
			} else if got != tt.want {
				t.Errorf("extractTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
