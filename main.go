package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	linkReplacement             = regexp.MustCompile(`(.?)<([^|>]+)\|([^>]+)>(.?)`) // refactor links
	appendLeadingEscape         = regexp.MustCompile(`(?m)^(#+\s|-\s)`)             // add escape for leading # or -
	backSlashReplacementFromGfm = regexp.MustCompile("\\\\([-*_`~#[]\\\\])")
	boldReplacement             = regexp.MustCompile("([\\s~_`]|^*)\\\\\\*([^\\s\\\\]+|[^*\\s\n]+[^*\n]*[^*\\s\n]+)\\\\\\*([\\s~_`\n]|$)") // find actual bold styles
	underscoreReplacement       = regexp.MustCompile("([*\\s~`]|^*)\\\\(_)([^\\s\n_]+|[^\\s\n_]+[^\n_]*[^\\s\n_]+)\\\\\\_([*\\s~`\n]|$)")  // find actual italic blocks
	tildeReplacement            = regexp.MustCompile("([*\\s_`]|^*)\\\\(~)([^\\s\n~]+|[^\\s\n~]+[^\n~]*[^\\s\n~]+)\\\\\\~([*\\s_`\n]|$)")  // find actual strike through blocks
	backtickReplacement         = regexp.MustCompile("\\\\`([^\n\\`]+)\\\\`")                                                              // find code
	parenthesisReplace          = regexp.MustCompile(`(?m)^\s*(\d+)\)`)                                                                    // find ordered list with format 1), 2) as UI render these as 1., 2.
	backslaskMultipler          = regexp.MustCompile(`(\\+)`)                                                                              // make \ double to make it slash symbol
	tildeToDoubleReplacement    = regexp.MustCompile(`([^\\]|^)~`)
	codeBlockReplacement        = regexp.MustCompile("(```)(.*)(```)")
	blockquoteReplacement       = regexp.MustCompile("(&gt; )([^\n]*)($)")
)

const tempStr string = "REPLACERSTR"

func SlackMarkdownToGeneral(slackMarkdown string) string {
	markdown := linkReplacement.ReplaceAllStringFunc(slackMarkdown, func(match string) string {
		submatches := linkReplacement.FindStringSubmatch(match)
		if len(submatches) > 4 {
			if submatches[1] == "`" && submatches[4] == "`" {
				return fmt.Sprintf("[`%s`](%s)", submatches[3], submatches[2])
			}
			return fmt.Sprintf("%s[%s](%s)%s", submatches[1], submatches[3], submatches[2], submatches[4])
		}
		return match
	})
	markdown = backslaskMultipler.ReplaceAllString(markdown, "$1$1")
	markdown = backSlashReplacementFromGfm.ReplaceAllString(markdown, "$1")

	markdown = strings.ReplaceAll(markdown, "*", "\\*")
	markdown = boldReplacement.ReplaceAllString(markdown, "$1**$2**$3")

	markdown = strings.ReplaceAll(markdown, "_", "\\_")
	markdown = underscoreReplacement.ReplaceAllString(markdown, "$1$2$3$2$4")

	markdown = strings.ReplaceAll(markdown, "~", "\\~")
	markdown = tildeReplacement.ReplaceAllString(markdown, "$1$2$3$2$4")

	markdown = codeBlockReplacement.ReplaceAllStringFunc(markdown, func(match string) string {
		submatches := codeBlockReplacement.FindStringSubmatch(match)
		return fmt.Sprintf("%s\n%s\n%s", tempStr, submatches[2], tempStr)
	})
	markdown = strings.ReplaceAll(markdown, "`", "\\`")

	markdown = blockquoteReplacement.ReplaceAllString(markdown, "$1$2$3\n")

	markdown = backtickReplacement.ReplaceAllString(markdown, "`$1`")
	markdown = strings.ReplaceAll(markdown, tempStr, "```")

	markdown = appendLeadingEscape.ReplaceAllString(markdown, "\\$1")
	markdown = parenthesisReplace.ReplaceAllString(markdown, "$1 )")

	return markdown
}

func MarkdownToHtmlMark(generalMarkdown string) string {
	markdown := tildeToDoubleReplacement.ReplaceAllString(generalMarkdown, "$1~~")
	markdown = strings.ReplaceAll(markdown, "\n", "\n\n") // double enter will render as <p> in html

	return markdown
}

func splitSentence(input string) []string {
	pattern := regexp.MustCompile(`(\*.*?\*|_.*?_|~.*?~|` + "`" + `.*?` + "`" + `|<.*?>|\s+|\b)`)
	matches := pattern.FindAllString(input, -1)

	var result []string
	var lastIndex int
	for _, match := range matches {
		index := pattern.FindStringIndex(input[lastIndex:])
		fmt.Println("")
		if beforeMatch := input[lastIndex : lastIndex+index[0]]; beforeMatch != "" {
			result = append(result, beforeMatch)
		}

		result = append(result, match)
		lastIndex += index[1] + index[0]
	}

	if lastIndex < len(input) {
		result = append(result, input[lastIndex:])
	}

	return result
}

// just add function, not used
func splitInput(input string) string {
	var splitedText []string

	catchIndex := 0
	for i, s := range input {
		if strings.Contains("*_~`<", string(s)) {
			catchIndex = i
		}
	}
}

func main() {
	input := "This is *complex* <https://github.com|style> for _~*my slack*~_ `message`"
	split := splitSentence(input)
	fmt.Printf("%q\n", split)
}
