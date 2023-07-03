// This file parses a given vm file and
// breaks up the words within a given line.
// This parser assumes input is always correct.
package main

import (
	"bufio"
	"os"
	"strings"
)

type VMProgram [][]string

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

func removeSpaces(vmLines []string) VMProgram {
	formattedVMLines := make([][]string, len(vmLines))
	for i := 0; i < len(vmLines); i++ {
		formattedVMLines[i] = strings.Split(vmLines[i], "")
	}
	return formattedVMLines
}

func formatVMLines(vmLines []string) VMProgram {
	vmLines = removeLineComments(vmLines)
	vmLines = removeEmptyLines(vmLines)
	return removeSpaces(vmLines)
}

func formatVMFile(file *os.File) VMProgram {
	vmLines := seperateEachLine(file)
	return formatVMLines(vmLines)
}

func parseVMFile(file *os.File) VMProgram {
	return formatVMFile(file)
}
