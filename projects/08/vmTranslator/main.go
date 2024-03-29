package main

import (
	"log"
	"os"
	"path/filepath"
)

const VM_EXTENSION = ".vm"

func isDir(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal("Could not open file path provided")
	}
	return fileInfo.IsDir()
}

func getFilesInDir(filePath string) []*os.File {
	files := []*os.File{}
	filepath.Walk(filePath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == VM_EXTENSION {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal("Error opening file in provided path")
				}
				files = append(files, file)
			}
		}
		return nil
	})
	return files
}

func getFile(filePath string) []*os.File {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error opening file in provided path")
	}
	return []*os.File{file}
}

func getFiles(filePath string) []*os.File {
	if isDir(filePath) {
		return getFilesInDir(filePath)
	}
	return getFile(filePath)
}

func runFile(file *os.File, cw *CodeWriter) {
	parser := getParser(file)
	for parser.hasMoreCommands() {
		parser.advance()
		if parser.commandType() == C_ARITHMETIC {
			cw.writeArithmetic(parser.arg1())
		} else if parser.commandType() == C_POP || parser.commandType() == C_PUSH {
			cw.writePushPop(parser.commandType(), parser.arg1(), parser.arg2())
		} else if parser.commandType() == C_LABEL {
			cw.writeLabel(parser.arg1())
		} else if parser.commandType() == C_GOTO {
			cw.writeGoto(parser.arg1())
		} else if parser.commandType() == C_IFGOTO {
			cw.writeIf(parser.arg1())
		} else if parser.commandType() == C_FUNCTION {
			cw.writeFunction(parser.arg1(), parser.arg2())
		} else if parser.commandType() == C_CALL {
			cw.writeCall(parser.arg1(), parser.arg2())
		} else if parser.commandType() == C_RETURN {
			cw.writeReturn()
		} else {
			log.Fatal("Invalid command found in vm file.")
		}
	}
}

func runFiles(filePath string) {
	files := getFiles(filePath)
	codeWriter := getCodeWriter(filePath)
	codeWriter.writeInit()
	for _, file := range files {
		codeWriter.setFileName(file.Name())
		runFile(file, codeWriter)
	}
	codeWriter.Close()
}

func main() {
	// Read input .vm file
	if len(os.Args) != 2 {
		log.Fatal("Incorrect number of inputs. Expected 2.")
	}
	filePath := os.Args[1]
	runFiles(filePath)
}
