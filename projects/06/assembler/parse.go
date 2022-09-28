package main

import "regexp"

type command string

const (
	AType  command = "atype"
	CType  command = "ctype"
	Ignore command = "ignore"
)

const (
	// Characters to remove from A & C type regex
	ATypeRegexToRemove   = "@"
	CTypeRegexToRemove   = "=|;"
	CommentRegexToRemove = "//[a-zA-Z0-9]*"
)

type parser struct{}

func (p *parser) isBlankLine(line string) bool {
	return len(line) == 0
}

func (p *parser) isLabel(line string) bool {
	return len(line) > 0 && line[0] == '('
}

func (p *parser) isComment(line string) bool {
	return len(line) > 1 && line[0] == '/' && line[1] == '/'
}

func (p *parser) isAType(line string) bool {
	return len(line) > 0 && line[0] == '@'
}

func (p *parser) isCType(line string) bool {
	return len(line) > 1 && !p.isComment(line) && !p.isLabel(line)
}

// Assumption that given files are always valid
func (p *parser) parseLabel(line string) string {
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

func (p *parser) parseAType(line string) []string {
	return splitStringFromRegex(removeComments(line), ATypeRegexToRemove)[1:]
}

func (p *parser) parseCType(line string) []string {
	return splitStringFromRegex(removeComments(line), CTypeRegexToRemove)
}

func (p *parser) parseLine(line string) (command, []string) {
	if p.isAType(line) {
		return AType, p.parseAType(line)
	} else if p.isCType(line) {
		return CType, p.parseCType(line)
	} else {
		return Ignore, nil
	}
}
