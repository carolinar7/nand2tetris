// This file parses a given vm file and
// breaks up the words within a given line.
// This parser assumes input is always correct.
package main

import (
	"bufio"
	"os"
)

type VMProgram struct {
	parsedVMLines [][]string
}

func seperateEachLine(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	vmLines := []string{}
	for scanner.Scan() {
		vmLines = append(vmLines, scanner.Text())
	}
	return vmLines
}

func parseVMFile(file *os.File) {
	seperateEachLine(file)
}
