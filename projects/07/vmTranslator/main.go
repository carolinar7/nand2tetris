package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const VM_EXTENSION = "vm"

func getFileNameAndTypeFromPath(filePath string) (string, string) {
	// Break down path
	filePathAsSlice := strings.Split(filePath, "/")
	fileStr := filePathAsSlice[len(filePathAsSlice)-1]
	// Break file by name and file type
	fileStrAsSlice := strings.Split(fileStr, ".")
	return fileStrAsSlice[0], fileStrAsSlice[1]
}

func main() {
	// Read input .vm file
	if len(os.Args) != 2 {
		log.Fatal("Incorrect number of inputs. Expected 2.")
	}
	filePath := os.Args[1]
	parser := getParser(filePath)
	for parser.hasMoreCommands() {
		parser.advance()
		fmt.Println(parser.commandType(), parser.arg1())
	}
	// outputASMFile, err := getASMFileFromPath(file.Name())
	// defer outputASMFile.Close()
	// writeToASMFile(outputASMFile, vmProgram)
}
