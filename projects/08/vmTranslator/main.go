package main

import (
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
	codeWriter := getCodeWriter(filePath)
	for parser.hasMoreCommands() {
		parser.advance()
		if parser.commandType() == C_ARITHMETIC {
			codeWriter.writeArithmetic(parser.arg1())
		} else if parser.commandType() == C_POP || parser.commandType() == C_PUSH {
			codeWriter.writePushPop(parser.commandType(), parser.arg1(), parser.arg2())
		}
	}
}
