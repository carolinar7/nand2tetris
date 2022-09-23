package main

type command string

const (
	AType  command = "atype"
	CType  command = "ctype"
	Ignore command = "ignore"
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

// TODO: Implement parsing fuctions
func (p *parser) parseAType(line string) []string { return nil }
func (p *parser) parseCType(line string) []string { return nil }

func (p *parser) parseLine(line string) (command, []string) {
	if p.isAType(line) {
		return AType, p.parseAType(line)
	} else if p.isCType(line) {
		return CType, p.parseCType(line)
	} else {
		return Ignore, nil
	}
}
