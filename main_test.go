package main

import (
	"reflect"
	"testing"
)

func TestSlackMarkdownToGeneral(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOutput string
	}{
		{
			name:       "Simple link",
			input:      "Visit <https://github.com|GitHub>",
			wantOutput: "Visit [GitHub](https://github.com)",
		},
		{
			name:       "Bold text",
			input:      "This is *bold* text",
			wantOutput: "This is **bold** text",
		},
		{
			name:       "Italic text",
			input:      "This is _italic_ text",
			wantOutput: "This is _italic_ text",
		},
		{
			name:       "Strikethrough text",
			input:      "This is ~strikethrough~ text",
			wantOutput: "This is ~strikethrough~ text",
		},
		{
			name:       "Code block",
			input:      "Here is `code`",
			wantOutput: "Here is `code`",
		},
		{
			name:       "Complex text",
			input:      "This is *complex* <https://github.com|style> for _~*my slack*~_ `message`",
			wantOutput: "This is **complex** [style](https://github.com) for _~**my slack**~_ `message`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlackMarkdownToGeneral(tt.input); got != tt.wantOutput {
				t.Errorf("SlackMarkdownToGeneral() = %v, want %v", got, tt.wantOutput)
			}
		})
	}
}

func TestMarkdownToHtmlMark(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOutput string
	}{
		{
			name:       "Simple markdown",
			input:      "This is *bold* text",
			wantOutput: "This is *bold* text",
		},
		{
			name:       "Tilde replacement",
			input:      "This is ~strikethrough~ text",
			wantOutput: "This is ~~strikethrough~~ text",
		},
		{
			name:       "Double enter for paragraph",
			input:      "This is line one\nThis is line two",
			wantOutput: "This is line one\n\nThis is line two",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarkdownToHtmlMark(tt.input); got != tt.wantOutput {
				t.Errorf("MarkdownToHtmlMark() = %v, want %v", got, tt.wantOutput)
			}
		})
	}
}

func TestSplitSentence(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOutput []string
	}{
		{
			name:       "Simple sentence",
			input:      "This is a test",
			wantOutput: []string{"This", " ", "is", " ", "a", " ", "test"},
		},
		{
			name:       "Sentence with formatting",
			input:      "This is *bold* and _italic_ text",
			wantOutput: []string{"This", " ", "is", " ", "*bold*", " ", "and", " ", "_italic_", " ", "text"},
		},
		{
			name:       "Complex sentence",
			input:      "This is *complex* <https://github.com|style> for _~*my slack*~_ `message`",
			wantOutput: []string{"This", " ", "is", " ", "*complex*", " ", "<https://github.com|style>", " ", "for", " ", "_~*my slack*~_", " ", "`message`"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitSentence(tt.input); !reflect.DeepEqual(got, tt.wantOutput) {
				t.Errorf("splitSentence() = %v, want %v", got, tt.wantOutput)
			}
		})
	}
}
