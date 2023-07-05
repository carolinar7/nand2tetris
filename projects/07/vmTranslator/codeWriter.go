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
// increment num clojure
func getIncrementNum() func() int {
	sum := -1
	return func() int {
		sum += 1
		return sum
	}
}

func decrementStackPointer(op []string) {
	// SP--
	op = append(op, "@SP")
	op = append(op, "AM=M-1")
}

func incrementStackPointer(op []string) {
	// SP++
	op = append(op, "@SP")
	op = append(op, "M=M+1")
}

func setComparison(op []string, comp string) {
	getNum := getIncrementNum()
	op = append(op, "M=M-D")
	// Jump to label if conditional is true
	op = append(op, fmt.Sprintf("@j%d", getNum()))
	op = append(op, fmt.Sprintf("M;%s", comp))
	// false
	op = append(op, "M=0")
	op = append(op, fmt.Sprintf("@j%dend", getNum()))
	op = append(op, "0;JMP")
	// true
	op = append(op, fmt.Sprintf("(j%d)", getNum()))
	op = append(op, "M=-1")
	op = append(op, fmt.Sprintf("(j%dend)", getNum()))
}

func getAdd() string {
	add := []string{}
	decrementStackPointer(add)
	// D=*SP
	add = append(add, "D=M")
	decrementStackPointer(add)
	// *SP=*SP+D
	add = append(add, "M=D+M")
	incrementStackPointer(add)
	return strings.Join(add, "\n")
}

func getSub() string {
	sub := []string{}
	decrementStackPointer(sub)
	// D=*SP
	sub = append(sub, "D=M")
	decrementStackPointer(sub)
	// *SP=*SP-D
	sub = append(sub, "M=M-D")
	incrementStackPointer(sub)
	return strings.Join(sub, "\n")
}

func getNeg() string {
	neg := []string{}
	decrementStackPointer(neg)
	// -*SP
	neg = append(neg, "M=-M")
	incrementStackPointer(neg)
	return strings.Join(neg, "\n")
}

func getEq() string {
	eq := []string{}
	decrementStackPointer(eq)
	// D=*SP
	eq = append(eq, "D=M")
	decrementStackPointer(eq)
	setComparison(eq, "JEQ")
	incrementStackPointer(eq)
	return strings.Join(eq, "\n")
}

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
		codeWriter.writeStringToOutput(getEq())
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
