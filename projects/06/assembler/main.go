package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getFileNameAndTypeFromPath(filePath string) (string, string) {
	// Break down path
	filePathAsSlice := strings.Split(filePath, "/")
	fileStr := filePathAsSlice[len(filePathAsSlice)-1]
	// Break file by name and file type
	fileStrAsSlice := strings.Split(fileStr, ".")
	return fileStrAsSlice[0], fileStrAsSlice[1]
}

func main() {
	// Read input .asm file
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Could not read file %v.", err)
	}
	defer file.Close()
	// Create output .hack file
	fileName, fileType := getFileNameAndTypeFromPath(file.Name())
	if fileType != "asm" {
		log.Fatal("File provided is not an .asm file type.")
	}
	fileOutName := fileName + ".hack"
	fileOut, err := os.Create(fileOutName)
	if err != nil {
		log.Fatalf("Could not make output file: %v", fileOutName)
	}
	defer fileOut.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// TODO: Parse each line
		if _, err := fileOut.WriteString(scanner.Text() + "\n"); err != nil {
			fmt.Print("Issue writing to output file.")
		}
	}
}
