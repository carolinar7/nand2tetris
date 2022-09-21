package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const ASSEMBLY_EXTENSION = "asm"

type program struct {
	lines []string
}

func (pg *program) removeSpaces() {
	for i := 0; i < len(pg.lines); i++ {
		pg.lines[i] = strings.Replace(pg.lines[i], " ", "", -1)
	}
}

func getFileNameAndTypeFromPath(filePath string) (string, string) {
	// Break down path
	filePathAsSlice := strings.Split(filePath, "/")
	fileStr := filePathAsSlice[len(filePathAsSlice)-1]
	// Break file by name and file type
	fileStrAsSlice := strings.Split(fileStr, ".")
	return fileStrAsSlice[0], fileStrAsSlice[1]
}

func getHackFileName(name string) string {
	return name + ".hack"
}

func getOutputHackFileFromPath(filePath string) (*os.File, error) {
	fileName, fileType := getFileNameAndTypeFromPath(filePath)
	if fileType != ASSEMBLY_EXTENSION {
		log.Fatal("File provided is not an .asm file type.")
	}
	fileOutName := getHackFileName(fileName)
	fileOut, err := os.Create(fileOutName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not make output file: %v", fileOutName))
	}
	defer fileOut.Close()
	return fileOut, nil
}

func main() {
	// Read input .asm file
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Could not read file %v.", err)
	}
	defer file.Close()
	fileOut, err := getOutputHackFileFromPath(file.Name())
	// TODO: REMOVE
	fmt.Print(fileOut)
	scanner := bufio.NewScanner(file)
	// We want to capture each line individually from scanner
	application := &program{[]string{}}
	for scanner.Scan() {
		application.lines = append(application.lines, scanner.Text())
	}
	application.removeSpaces()
}
