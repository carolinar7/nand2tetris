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

func getASMFileName(path, fileName string) string {
	return path + "/" + fileName + "." + ASM_EXTENSION
}

func getOutputFileFromInputPath(filePath string) (*os.File, error) {
	fileName, _, path := getFileNameAndTypeFromPath(filePath)
	fileOutName := getASMFileName(path, fileName)
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

func joinStrings(strs []string) string {
	str := strings.Join(strs, "\n")
	return fmt.Sprintf("%s\n", str)
}

// Arithmetic / Logical Commands
func decrementStackPointer(op []string) []string {
	// SP--
	op = append(op, "@SP")
	op = append(op, "AM=M-1")
	return op
}

func incrementStackPointer(op []string) []string {
	// SP++
	op = append(op, "@SP")
	op = append(op, "M=M+1")
	return op
}

var num int = 0

func setComparison(op []string, comp string) []string {
	op = append(op, "D=M-D")
	// Jump to label if conditional is true
	op = append(op, fmt.Sprintf("@j%d", num))
	op = append(op, fmt.Sprintf("D;%s", comp))
	// false
	op = append(op, "@SP")
	op = append(op, "A=M")
	op = append(op, "M=0")
	op = append(op, fmt.Sprintf("@j%dend", num))
	op = append(op, "0;JMP")
	// true
	op = append(op, fmt.Sprintf("(j%d)", num))
	op = append(op, "@SP")
	op = append(op, "A=M")
	op = append(op, "M=-1")
	op = append(op, fmt.Sprintf("(j%dend)", num))
	num++
	return op
}

func getAdd() string {
	add := []string{}
	add = decrementStackPointer(add[:])
	// D=*SP
	add = append(add, "D=M")
	add = decrementStackPointer(add)
	// *SP=*SP+D
	add = append(add, "M=D+M")
	add = incrementStackPointer(add[:])
	return joinStrings(add)
}

func getSub() string {
	sub := []string{}
	sub = decrementStackPointer(sub)
	// D=*SP
	sub = append(sub, "D=M")
	sub = decrementStackPointer(sub)
	// *SP=*SP-D
	sub = append(sub, "M=M-D")
	sub = incrementStackPointer(sub)
	return joinStrings(sub)
}

func getNeg() string {
	neg := []string{}
	neg = decrementStackPointer(neg)
	// -*SP
	neg = append(neg, "M=-M")
	neg = incrementStackPointer(neg)
	return joinStrings(neg)
}

func getEq() string {
	eq := []string{}
	eq = decrementStackPointer(eq)
	// D=*SP
	eq = append(eq, "D=M")
	eq = decrementStackPointer(eq)
	eq = setComparison(eq, "JEQ")
	eq = incrementStackPointer(eq)
	return joinStrings(eq)
}

func getGt() string {
	gt := []string{}
	gt = decrementStackPointer(gt)
	// D=*SP
	gt = append(gt, "D=M")
	gt = decrementStackPointer(gt)
	gt = setComparison(gt, "JGT")
	gt = incrementStackPointer(gt)
	return joinStrings(gt)
}

func getLt() string {
	lt := []string{}
	lt = decrementStackPointer(lt)
	// D=*SP
	lt = append(lt, "D=M")
	lt = decrementStackPointer(lt)
	lt = setComparison(lt, "JLT")
	lt = incrementStackPointer(lt)
	return joinStrings(lt)
}

func getAnd() string {
	and := []string{}
	and = decrementStackPointer(and)
	// D=*SP
	and = append(and, "D=M")
	and = decrementStackPointer(and)
	// M=M&D
	and = append(and, "M=D&M")
	and = incrementStackPointer(and)
	return joinStrings(and)
}

func getOr() string {
	or := []string{}
	or = decrementStackPointer(or)
	// D=*SP
	or = append(or, "D=M")
	or = decrementStackPointer(or)
	// M=M&D
	or = append(or, "M=D|M")
	or = incrementStackPointer(or)
	return joinStrings(or)
}

func getNot() string {
	not := []string{}
	not = decrementStackPointer(not)
	// M=!M
	not = append(not, "M=!M")
	not = incrementStackPointer(not)
	return joinStrings(not)
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
const LOCAL_ABBR = "LCL"
const ARGUMENT_ABBR = "ARG"
const THIS_ABBR = "THIS"
const THAT_ABBR = "THAT"

func segmentPointersPush(segment string, idx int) string {
	push := []string{}
	// D=i
	push = append(push, fmt.Sprintf("@%d", idx))
	push = append(push, "D=A")
	// addr=Seg+i
	push = append(push, fmt.Sprintf("@%s", segment))
	push = append(push, "A=D+M")
	// D=*addr
	push = append(push, "D=M")
	// *SP=D
	push = append(push, "@SP")
	push = append(push, "A=M")
	push = append(push, "M=D")
	push = incrementStackPointer(push)
	return joinStrings(push)
}

func getLocalPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(LOCAL_ABBR, idx)
	}
	return ""
}

func getArgumentPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(ARGUMENT_ABBR, idx)
	}
	return ""
}

func getThisPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(THIS_ABBR, idx)
	}
	return ""
}

func getThatPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(THAT_ABBR, idx)
	}
	return ""
}

func getConstantPush(idx int) string {
	constant := []string{}
	// D=i
	constant = append(constant, fmt.Sprintf("@%d", idx))
	constant = append(constant, "D=A")
	// *SP=D
	constant = append(constant, "@SP")
	constant = append(constant, "A=M")
	constant = append(constant, "M=D")
	constant = incrementStackPointer(constant)
	return joinStrings(constant)
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
