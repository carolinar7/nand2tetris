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
	"strings"
)

const ASM_EXTENSION = "asm"

type CodeWriter struct {
	outputFile *os.File
}

func (codeWriter *CodeWriter) writeStringToOutput(str string) {
	codeWriter.outputFile.WriteString(str)
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
func getAdd() string {
	add := []string{}
	// D=*SP
	add = append(add, "@SP")
	add = append(add, "A=M")
	add = append(add, "D=M")
	// SP--
	add = append(add, "@SP")
	add = append(add, "M=M-1")
	// *SP=*SP+D
	add = append(add, "@SP")
	add = append(add, "A=M")
	add = append(add, "M=D+M")
	return strings.Join(add, "\n")
}

func getSub() string {
	sub := []string{}
	// D=*SP
	sub = append(sub, "@SP")
	sub = append(sub, "A=M")
	sub = append(sub, "D=M")
	// SP--
	sub = append(sub, "@SP")
	sub = append(sub, "M=M-1")
	// *SP=*SP-D
	sub = append(sub, "@SP")
	sub = append(sub, "A=M")
	sub = append(sub, "M=M-D")
	return strings.Join(sub, "\n")
}

func getNeg() string {
	neg := []string{}
	// -*SP
	neg = append(neg, "@SP")
	neg = append(neg, "A=M")
	neg = append(neg, "M=-M")
	return strings.Join(neg, "\n")
}

// eq

// gt

// lt

// and

// or

// not

func (codeWriter *CodeWriter) writeArithmetic(command string) {
	switch command {
	case "add":
		codeWriter.writeStringToOutput(getAdd())
	case "sub":
		codeWriter.writeStringToOutput(getSub())
	case "neg":
		codeWriter.writeStringToOutput(getNeg())
	case "eq":
	case "gt":
	case "lt":
	case "and":
	case "or":
	case "not":
	}
}

// Memory Access Commands
// local

// argument

// this

// that

// constant

// static

// temp

// pointer
