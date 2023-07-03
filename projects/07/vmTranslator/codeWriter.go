// This file is in charge of code generation or translating
// the parsed vm code to its respective asm code. It is
// broken down into two sections: arithmetic/logical commands
// and memory access commands.
//
// Arithmetic / Locagical commands contain code generation in
// asm for add, sub, neg, eq, gt, lt, and, or & not.
//
// Memory access commands follow the format:
// 	pop segment i OR push segment i
// Where a segment can either represent: local, argument,
// this, that, constant, static, temp, & pointer. These
// are segments within RAM. In addition, i represents a
// numerical value that determines the shift off the respective
// base pointer.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

const ASM_EXTENSION = "asm"

type CodeWriter struct {
	outputFile *os.File
}

func getASMFileName(fileName string) string {
	return fileName + "." + ASM_EXTENSION
}

func getOutputFileFromInputPath(filePath string) (*os.File, error) {
	fileName, _ := getFileNameAndTypeFromPath(filePath)
	fileOutName := getASMFileName(fileName)
	fileOut, err := os.Create(fileOutName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not make output file: %v.", fileOutName))
	}
	return fileOut, nil
}

func getCodeWriter(filePath string) *CodeWriter {
	outputfile, err := getOutputFileFromInputPath(filePath)
	if err != nil {
		log.Fatal("Could not create outputfile")
	}
	return &CodeWriter{outputFile: outputfile}
}

// Arithmetic / Logical Commands
// add

// sub

// neg

// eq

// gt

// lt

// and

// or

// not

// Memory Access Commands
// local

// argument

// this

// that

// constant

// static

// temp

// pointer
