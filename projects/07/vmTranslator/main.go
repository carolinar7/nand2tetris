package main

import (
	"fmt"
	"log"
	"os"
)

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
	codeWriter := getCodeWriter(filePath)
	codeWriter.outputFile.WriteString("Hello World")
}
