// This file parses a given vm file and
// breaks up the words within a given line.
// This parser assumes input is always correct.
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	// iota auto-increments from 0
	C_ARITHMETIC int = iota
	C_PUSH
	C_POP
)
const VM_EXTENSION = "vm"

type VMProgram []string
type Command []string

var CommandType = map[string]int{
	"add":  C_ARITHMETIC,
	"sub":  C_ARITHMETIC,
	"neg":  C_ARITHMETIC,
	"eq":   C_ARITHMETIC,
	"gt":   C_ARITHMETIC,
	"lt":   C_ARITHMETIC,
	"and":  C_ARITHMETIC,
	"or":   C_ARITHMETIC,
	"not":  C_ARITHMETIC,
	"push": C_PUSH,
	"pop":  C_POP,
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

func removeLineComments(vmLines []string) []string {
	newLines := []string{}
	for _, line := range vmLines {
		if !isComment(line) {
			newLines = append(newLines, line)
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

func getFileNameAndTypeFromPath(filePath string) (string, string, string) {
	// Break down path
	filePathAsSlice := strings.Split(filePath, "/")
	fileStr := filePathAsSlice[len(filePathAsSlice)-1]
	// Break file by name and file type
	fileStrAsSlice := strings.Split(fileStr, ".")
	return fileStrAsSlice[0], fileStrAsSlice[1], strings.Join(filePathAsSlice[:len(filePathAsSlice)-1], "/")
}

func openVMFile(filePath string) *os.File {
	_, fileType, _ := getFileNameAndTypeFromPath(filePath)
	if fileType != VM_EXTENSION {
		log.Fatal("File Provided is not a .vm file type.")
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Could not read file %v.", err)
	}
	return file
}

func getParser(filePath string) *Parser {
	file := openVMFile(filePath)
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
	if CommandType[parser.command[0]] == C_ARITHMETIC {
		return parser.command[0]
	}
	return parser.command[1]
}

func (parser *Parser) arg2() int {
	if parser.command == nil || len(parser.command) < 2 {
		log.Fatal("Attempted to access arg2 that does not exist")
	}
	arg2, err := strconv.Atoi(parser.command[2])
	if err != nil {
		log.Fatalf("Arg2 is not valid: %v", err)
	}
	return arg2
}
