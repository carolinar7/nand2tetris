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

const ASM_EXTENSION = ".asm"

type CodeWriter struct {
	outputFile   *os.File
	incrementNum int
	fileName     string
}

func (codeWriter *CodeWriter) writeStringToOutput(str string) {
	codeWriter.outputFile.WriteString(str)
}

func getFileNameAndTypeFromPath(filePath string) (string, string, string) {
	// Break down path
	filePathAsSlice := strings.Split(filePath, "\\")
	fileStr := filePathAsSlice[len(filePathAsSlice)-1]
	// Break file by name and file type
	fileStrAsSlice := strings.Split(fileStr, ".")
	fmt.Println(filePathAsSlice)
	return fileStrAsSlice[0], fileStrAsSlice[1], strings.Join(filePathAsSlice[:len(filePathAsSlice)-1], "/")
}

func getASMFileName(path, fileName string) string {
	fmt.Println(path, "hey", fileName)
	return path + "/" + fileName + ASM_EXTENSION
}

func getOutputFileFromInputPath(filePath string) (*os.File, string, error) {
	fileName, _, path := getFileNameAndTypeFromPath(filePath)
	fileOutName := getASMFileName(path, fileName)
	fileOut, err := os.Create(fileOutName)
	if err != nil {
		return nil, "", errors.New(fmt.Sprintf("Could not make output file: %v", fileOutName))
	}
	return fileOut, fileName, nil
}

func getCodeWriter(file *os.File) *CodeWriter {
	outputfile, fileName, err := getOutputFileFromInputPath(file.Name())
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &CodeWriter{outputFile: outputfile, incrementNum: 0, fileName: fileName}
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

func (cw *CodeWriter) setComparison(op []string, comp string) []string {
	op = append(op, "D=M-D")
	// Jump to label if conditional is true
	op = append(op, fmt.Sprintf("@j%d", cw.incrementNum))
	op = append(op, fmt.Sprintf("D;%s", comp))
	// false
	op = append(op, "@SP")
	op = append(op, "A=M")
	op = append(op, "M=0")
	op = append(op, fmt.Sprintf("@j%dend", cw.incrementNum))
	op = append(op, "0;JMP")
	// true
	op = append(op, fmt.Sprintf("(j%d)", cw.incrementNum))
	op = append(op, "@SP")
	op = append(op, "A=M")
	op = append(op, "M=-1")
	op = append(op, fmt.Sprintf("(j%dend)", cw.incrementNum))
	cw.incrementNum++
	return op
}

func (cw *CodeWriter) getAdd() string {
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

func (cw *CodeWriter) getSub() string {
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

func (cw *CodeWriter) getNeg() string {
	neg := []string{}
	neg = decrementStackPointer(neg)
	// -*SP
	neg = append(neg, "M=-M")
	neg = incrementStackPointer(neg)
	return joinStrings(neg)
}

func (cw *CodeWriter) getEq() string {
	eq := []string{}
	eq = decrementStackPointer(eq)
	// D=*SP
	eq = append(eq, "D=M")
	eq = decrementStackPointer(eq)
	eq = cw.setComparison(eq, "JEQ")
	eq = incrementStackPointer(eq)
	return joinStrings(eq)
}

func (cw *CodeWriter) getGt() string {
	gt := []string{}
	gt = decrementStackPointer(gt)
	// D=*SP
	gt = append(gt, "D=M")
	gt = decrementStackPointer(gt)
	gt = cw.setComparison(gt, "JGT")
	gt = incrementStackPointer(gt)
	return joinStrings(gt)
}

func (cw *CodeWriter) getLt() string {
	lt := []string{}
	lt = decrementStackPointer(lt)
	// D=*SP
	lt = append(lt, "D=M")
	lt = decrementStackPointer(lt)
	lt = cw.setComparison(lt, "JLT")
	lt = incrementStackPointer(lt)
	return joinStrings(lt)
}

func (cw *CodeWriter) getAnd() string {
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

func (cw *CodeWriter) getOr() string {
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

func (cw *CodeWriter) getNot() string {
	not := []string{}
	not = decrementStackPointer(not)
	// M=!M
	not = append(not, "M=!M")
	not = incrementStackPointer(not)
	return joinStrings(not)
}

func (cw *CodeWriter) writeArithmetic(command string) {
	switch command {
	case "add":
		cw.writeStringToOutput(cw.getAdd())
	case "sub":
		cw.writeStringToOutput(cw.getSub())
	case "neg":
		cw.writeStringToOutput(cw.getNeg())
	case "eq":
		cw.writeStringToOutput(cw.getEq())
	case "gt":
		cw.writeStringToOutput(cw.getGt())
	case "lt":
		cw.writeStringToOutput(cw.getLt())
	case "and":
		cw.writeStringToOutput(cw.getAnd())
	case "or":
		cw.writeStringToOutput(cw.getOr())
	case "not":
		cw.writeStringToOutput(cw.getNot())
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

func segmentPointerPop(segment string, idx int) string {
	pop := []string{}
	// D=i
	pop = append(pop, fmt.Sprintf("@%d", idx))
	pop = append(pop, "D=A")
	// addr=Seg+i
	pop = append(pop, fmt.Sprintf("@%s", segment))
	pop = append(pop, "D=D+M")
	pop = append(pop, "@R13")
	pop = append(pop, "M=D")
	pop = decrementStackPointer(pop)
	// *addr=*SP
	pop = append(pop, "D=M")
	pop = append(pop, "@R13")
	pop = append(pop, "A=M")
	pop = append(pop, "M=D")
	return joinStrings(pop)
}

func (cw *CodeWriter) getLocalPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(LOCAL_ABBR, idx)
	}
	return segmentPointerPop(LOCAL_ABBR, idx)
}

func (cw *CodeWriter) getArgumentPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(ARGUMENT_ABBR, idx)
	}
	return segmentPointerPop(ARGUMENT_ABBR, idx)
}

func (cw *CodeWriter) getThisPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(THIS_ABBR, idx)
	}
	return segmentPointerPop(THIS_ABBR, idx)
}

func (cw *CodeWriter) getThatPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(THAT_ABBR, idx)
	}
	return segmentPointerPop(THAT_ABBR, idx)
}

func (cw *CodeWriter) getConstantPush(idx int) string {
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

func pushStatic(idx int, fileName string) string {
	push := []string{}
	// D=fileName.i
	push = append(push, fmt.Sprintf("@%s.%d", fileName, idx))
	push = append(push, "D=M")
	// *SP=D
	push = append(push, "@SP")
	push = append(push, "A=M")
	push = append(push, "M=D")
	push = incrementStackPointer(push)
	return joinStrings(push)
}

func popStatic(idx int, fileName string) string {
	pop := []string{}
	pop = decrementStackPointer(pop)
	// D=*SP
	pop = append(pop, "D=M")
	// fileName.i=D
	pop = append(pop, fmt.Sprintf("@%s.%d", fileName, idx))
	pop = append(pop, "M=D")
	return joinStrings(pop)
}

func (cw *CodeWriter) getStaticPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return pushStatic(idx, cw.fileName)
	}
	return popStatic(idx, cw.fileName)
}

func pushTemp(idx int) string {
	push := []string{}
	// D=i
	push = append(push, fmt.Sprintf("@%d", idx))
	push = append(push, "D=A")
	// addr=Seg+i
	push = append(push, "@5")
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

func popTemp(idx int) string {
	pop := []string{}
	// D=i
	pop = append(pop, fmt.Sprintf("@%d", idx))
	pop = append(pop, "D=A")
	// addr=Seg+i
	pop = append(pop, "@5")
	pop = append(pop, "D=D+A")
	pop = append(pop, "@R13")
	pop = append(pop, "M=D")
	pop = decrementStackPointer(pop)
	// *addr=*SP
	pop = append(pop, "D=M")
	pop = append(pop, "@R13")
	pop = append(pop, "A=M")
	pop = append(pop, "M=D")
	return joinStrings(pop)
}

func (cw *CodeWriter) getTempPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return pushTemp(idx)
	}
	return popTemp(idx)
}

const THIS = 0
const THAT = 1

func pushPointer(idx int, pt string) string {
	push := []string{}
	// D=THIS/THAT
	push = append(push, fmt.Sprintf("@%s", pt))
	push = append(push, "D=M")
	// *SP=D
	push = append(push, "@SP")
	push = append(push, "A=M")
	push = append(push, "M=D")
	push = incrementStackPointer(push)
	return joinStrings(push)
}

func popPointer(idx int, pt string) string {
	pop := []string{}
	pop = decrementStackPointer(pop)
	// D=*SP
	pop = append(pop, "D=M")
	// THIS/THAT=D
	pop = append(pop, fmt.Sprintf("@%s", pt))
	pop = append(pop, "M=D")
	return joinStrings(pop)
}

func (cw *CodeWriter) getPointerPushPop(pushPop int, idx int) string {
	var pointerType string
	if idx == THIS {
		pointerType = "THIS"
	} else if idx == THAT {
		pointerType = "THAT"
	}
	if pushPop == C_PUSH {
		return pushPointer(idx, pointerType)
	}
	return popPointer(idx, pointerType)
}

func (cw *CodeWriter) writePushPop(pushPop int, segment string, idx int) {
	switch segment {
	case "local":
		cw.writeStringToOutput(cw.getLocalPushPop(pushPop, idx))
	case "argument":
		cw.writeStringToOutput(cw.getArgumentPushPop(pushPop, idx))
	case "this":
		cw.writeStringToOutput(cw.getThisPushPop(pushPop, idx))
	case "that":
		cw.writeStringToOutput(cw.getThatPushPop(pushPop, idx))
	case "constant":
		cw.writeStringToOutput(cw.getConstantPush(idx))
	case "static":
		cw.writeStringToOutput(cw.getStaticPushPop(pushPop, idx))
	case "temp":
		cw.writeStringToOutput(cw.getTempPushPop(pushPop, idx))
	case "pointer":
		cw.writeStringToOutput(cw.getPointerPushPop(pushPop, idx))
	}
}
