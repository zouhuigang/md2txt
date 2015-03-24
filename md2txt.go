package md2txt

import (
	"regexp"
)

const (
	BASIC = iota // Basic Markdown
	GFM          // Github Flavored Markdown
)

// Md2TxtString converts markdown format input string
// to pure text string.
func Md2TxtString(input string) string {
	var output string
	// remove #+ at the beginning.
	output = regexp.MustCompile("(?m)^#+").ReplaceAllString(input, "")
	output = regexp.MustCompile("(?m)^(>)+").ReplaceAllString(:wq
	return output
}
