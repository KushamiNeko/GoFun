package web

import (
	"regexp"
)

var (
	leadingSpace *regexp.Regexp

	htmlLine   *regexp.Regexp
	jsCSSLine  *regexp.Regexp
	jsCSSBlock *regexp.Regexp
)

func init() {
	leadingSpace = regexp.MustCompile(`(?m)^\s+`)
	htmlLine = regexp.MustCompile(`\s*<!--\s*(?:[\s\S]*?)\s*-->`)
	jsCSSLine = regexp.MustCompile(`(?:[^:"']|^)\/\/\s*(?:[^"'\n]+)`)
	jsCSSBlock = regexp.MustCompile(`\s*\/\*\s*[\s\S]*?\s*\*\/`)
}

func CleanAll(content []byte) []byte {
	tar := []byte("")

	content = jsCSSLine.ReplaceAll(content, tar)
	content = jsCSSBlock.ReplaceAll(content, tar)
	content = htmlLine.ReplaceAll(content, tar)
	content = leadingSpace.ReplaceAll(content, tar)
	return content
}

func CleanHTML(content []byte) []byte {
	return htmlLine.ReplaceAll(content, []byte(""))
}

func CleanJS(content []byte) []byte {
	tar := []byte("")

	content = jsCSSLine.ReplaceAll(content, tar)
	content = jsCSSBlock.ReplaceAll(content, tar)
	return content
}

func CleanLeadingSpace(content []byte) []byte {
	return leadingSpace.ReplaceAll(content, []byte(""))
}
