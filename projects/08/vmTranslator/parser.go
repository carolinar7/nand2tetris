// This file parses a given vm file and
// breaks up the words within a given line.
// This parser assumes input is always correct.
package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	// iota auto-increments from 0
	C_ARITHMETIC int = iota
	C_PUSH
	C_POP
	C_GOTO
	C_IFGOTO
	C_LABEL
	C_CALL
	C_FUNCTION
	C_RETURN
)

const CommentRegexToRemove = "//*"

type VMProgram []string
type Command []string

var CommandType = map[string]int{
	"add":      C_ARITHMETIC,
	"sub":      C_ARITHMETIC,
	"neg":      C_ARITHMETIC,
	"eq":       C_ARITHMETIC,
	"gt":       C_ARITHMETIC,
	"lt":       C_ARITHMETIC,
	"and":      C_ARITHMETIC,
	"or":       C_ARITHMETIC,
	"not":      C_ARITHMETIC,
	"push":     C_PUSH,
	"pop":      C_POP,
	"goto":     C_GOTO,
	"if-goto":  C_IFGOTO,
	"label":    C_LABEL,
	"call":     C_CALL,
	"function": C_FUNCTION,
	"return":   C_RETURN,
}

type Parser struct {
	vmLines  VMProgram
	lineNum  int
	command  Command
	numLines int
}

func seperateEachLine(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	vmLines := []string{}
	for scanner.Scan() {
		vmLines = append(vmLines, scanner.Text())
	}
	return vmLines
}

func isComment(line string) bool {
	return len(line) > 1 && line[0] == '/' && line[1] == '/'
}

func removeInlineComments(line string) string {
	re := regexp.MustCompile(CommentRegexToRemove)
	return re.Split(line, -1)[0]
}

func removeLineComments(vmLines []string) []string {
	newLines := []string{}
	for _, line := range vmLines {
		if !isComment(line) {
			newLines = append(newLines, removeInlineComments(line))
		}
	}
	vmLines = newLines
	return vmLines
}

func isBlankLine(line string) bool {
	return len(line) == 0
}

func removeEmptyLines(vmLines []string) []string {
	newLines := []string{}
	for _, line := range vmLines {
		if !isBlankLine(line) {
			newLines = append(newLines, line)
		}
	}
	vmLines = newLines
	return vmLines
}

func formatVMLines(vmLines []string) VMProgram {
	vmLines = removeLineComments(vmLines)
	return removeEmptyLines(vmLines)
}

func formatVMFile(file *os.File) VMProgram {
	vmLines := seperateEachLine(file)
	return formatVMLines(vmLines)
}

func parseVMFile(file *os.File) VMProgram {
	return formatVMFile(file)
}

func getParser(file *os.File) *Parser {
	vmProgram := parseVMFile(file)
	defer file.Close()
	// lineNum is set to -1 because initially there is no command set.
	// When the parser advances it will move to 0 and read the command,
	// and so forth.
	return &Parser{vmLines: vmProgram, lineNum: -1, numLines: len(vmProgram)}
}

func (parser *Parser) hasMoreCommands() bool {
	return parser.lineNum < parser.numLines-1
}

func tokenizeCommand(command string) Command {
	return strings.Split(command, " ")
}

func (parser *Parser) advance() {
	parser.lineNum = parser.lineNum + 1
	parser.command = tokenizeCommand(parser.vmLines[parser.lineNum])
}

func (parser *Parser) commandType() int {
	if parser.command == nil || len(parser.command) == 0 {
		log.Fatal("Attempted to access uninstantiated command")
	}
	return CommandType[parser.command[0]]
}

func (parser *Parser) arg1() string {
	if parser.command == nil || len(parser.command) == 0 {
		log.Fatal("Attempted to access uninstantiated command")
	}
	if CommandType[parser.command[0]] == C_ARITHMETIC || CommandType[parser.command[0]] == C_RETURN {
		return parser.command[0]
	}
	return parser.command[1]
}

func (parser *Parser) arg2() int {
	if parser.command == nil || len(parser.command) < 3 {
		log.Fatal("Attempted to access arg2 that does not exist")
	}
	arg2, err := strconv.Atoi(parser.command[2])
	if err != nil {
		log.Fatalf("Arg2 is not valid: %v", err)
	}
	return arg2
}
