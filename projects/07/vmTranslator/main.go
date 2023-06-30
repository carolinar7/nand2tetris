package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const VM_EXTENSION = "vm"
const ASM_EXTENSION = "asm"

func getFileNameAndTypeFromPath(filePath string) (string, string) {
	// Break down path
	filePathAsSlice := strings.Split(filePath, "/")
	fileStr := filePathAsSlice[len(filePathAsSlice)-1]
	// Break file by name and file type
	fileStrAsSlice := strings.Split(fileStr, ".")
	return fileStrAsSlice[0], fileStrAsSlice[1]
}

func getASMFileName(fileName string) string {
	return fileName + "." + ASM_EXTENSION
}

func getASMFileFromPath(filePath string) (*os.File, error) {
	fileName, fileType := getFileNameAndTypeFromPath(filePath)
	if fileType != VM_EXTENSION {
		log.Fatal("File Provided is not a .vm file type.")
	}
	fileOutName := getASMFileName(fileName)
	fileOut, err := os.Create(fileOutName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not make output file: %v.", fileOutName))
	}
	return fileOut, nil
}

func main() {
	// Read input .vm file
	if len(os.Args) != 2 {
		log.Fatal("Incorrect number of inputs. Expected 2.")
	}
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Could not read file %v.", err)
	}
	defer file.Close()
	outputASMFile, err := getASMFileFromPath(file.Name())
	defer outputASMFile.Close()
	parseVMFile(file)
}
