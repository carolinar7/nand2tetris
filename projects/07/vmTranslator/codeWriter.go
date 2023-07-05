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

func getGt() string {
	gt := []string{}
	decrementStackPointer(gt)
	// D=*SP
	gt = append(gt, "D=M")
	decrementStackPointer(gt)
	setComparison(gt, "JGT")
	incrementStackPointer(gt)
	return strings.Join(gt, "\n")
}

func getLt() string {
	lt := []string{}
	decrementStackPointer(lt)
	// D=*SP
	lt = append(lt, "D=M")
	decrementStackPointer(lt)
	setComparison(lt, "JLT")
	incrementStackPointer(lt)
	return strings.Join(lt, "\n")
}

func getAnd() string {
	and := []string{}
	decrementStackPointer(and)
	// D=*SP
	and = append(and, "D=M")
	decrementStackPointer(and)
	// M=M&D
	and = append(and, "M=D&M")
	incrementStackPointer(and)
	return strings.Join(and, "\n")
}

func getOr() string {
	or := []string{}
	decrementStackPointer(or)
	// D=*SP
	or = append(or, "D=M")
	decrementStackPointer(or)
	// M=M&D
	or = append(or, "M=D|M")
	incrementStackPointer(or)
	return strings.Join(or, "\n")
}

func getNot() string {
	not := []string{}
	decrementStackPointer(not)
	// M=!M
	not = append(not, "M=!M")
	incrementStackPointer(not)
	return strings.Join(not, "\n")
}

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
		codeWriter.writeStringToOutput(getGt())
	case "lt":
		codeWriter.writeStringToOutput(getLt())
	case "and":
		codeWriter.writeStringToOutput(getAnd())
	case "or":
		codeWriter.writeStringToOutput(getOr())
	case "not":
		codeWriter.writeStringToOutput(getNot())
	}
}

// Memory Access Commands
func getLocalPushPop(pushPop int, idx int) string {
	return ""
}

func getArgumentPushPop(pushPop int, idx int) string {
	return ""
}

func getThisPushPop(pushPop int, idx int) string {
	return ""
}

func getThatPushPop(pushPop int, idx int) string {
	return ""
}

func getConstantPush(idx int) string {
	constant := []string{}
	// D=i
	constant = append(constant, fmt.Sprintf("D=%d", idx))
	// *SP=D
	constant = append(constant, "@SP")
	constant = append(constant, "M=D")
	incrementStackPointer(constant)
	return strings.Join(constant, "\n")
}

func getStaticPushPop(pushPop int, idx int) string {
	return ""
}

func getTempPushPop(pushPop int, idx int) string {
	return ""
}

func getPointerPushPop(pushPop int, idx int) string {
	return ""
}

func (codeWriter *CodeWriter) writePushPop(pushPop int, segment string, idx int) {
	switch segment {
	case "local":
		codeWriter.writeStringToOutput(getLocalPushPop(pushPop, idx))
	case "argument":
		codeWriter.writeStringToOutput(getArgumentPushPop(pushPop, idx))
	case "this":
		codeWriter.writeStringToOutput(getThisPushPop(pushPop, idx))
	case "that":
		codeWriter.writeStringToOutput(getThatPushPop(pushPop, idx))
	case "constant":
		codeWriter.writeStringToOutput(getConstantPush(idx))
	case "static":
		codeWriter.writeStringToOutput(getStaticPushPop(pushPop, idx))
	case "temp":
		codeWriter.writeStringToOutput(getTempPushPop(pushPop, idx))
	case "pointer":
		codeWriter.writeStringToOutput(getPointerPushPop(pushPop, idx))
	}
}
