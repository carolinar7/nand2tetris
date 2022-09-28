package main

import "regexp"

type command string

const (
	AType command = "atype"
	CType command = "ctype"
)

const (
	// Characters to remove from A & C type regex
	ATypeRegexToRemove   = "@"
	CTypeRegexToRemove   = "=|;"
	CommentRegexToRemove = "//[a-zA-Z0-9]*"
)

type parser struct{}

func isBlankLine(line string) bool {
	return len(line) == 0
}

func isLabel(line string) bool {
	return len(line) > 0 && line[0] == '('
}

func isComment(line string) bool {
	return len(line) > 1 && line[0] == '/' && line[1] == '/'
}

func isAType(line string) bool {
	return len(line) > 0 && line[0] == '@'
}

func isCType(line string) bool {
	return len(line) > 1 && !isComment(line) && !isLabel(line)
}

// Assumption that given files are always valid
func parseLabel(line string) string {
	return line[1 : len(line)-1]
}

func splitStringFromRegex(line, regex string) []string {
	re := regexp.MustCompile(regex)
	return re.Split(line, -1)
}

func removeComments(line string) string {
	re := regexp.MustCompile(CommentRegexToRemove)
	return re.Split(line, -1)[0]
}

func parseAType(line string) []string {
	return splitStringFromRegex(removeComments(line), ATypeRegexToRemove)[1:]
}

func parseCType(line string) []string {
	return splitStringFromRegex(removeComments(line), CTypeRegexToRemove)
}

func parseLine(line string) *fieldTuple {
	if isAType(line) {
		return &fieldTuple{AType, parseAType(line)}
	} else if isCType(line) {
		return &fieldTuple{CType, parseCType(line)}
	}
	return nil
}
