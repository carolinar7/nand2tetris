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
const TRUE = -1
const FALSE = 0

type CodeWriter struct {
	outputFile *os.File
	num        int
	fileName   string
}

func (cw *CodeWriter) incrementNum() int {
	num := cw.num
	cw.num++
	return num
}

func (codeWriter *CodeWriter) writeStringToOutput(str string) {
	codeWriter.outputFile.WriteString(str)
}

func getFileNameAndDirectory(filePath string) (string, string) {
	// Break down path
	filePathAsSlice := strings.Split(filePath, "/")
	fileStr := filePathAsSlice[len(filePathAsSlice)-1]
	// Break file by name and file type
	fileStrAsSlice := strings.Split(fileStr, ".")
	if isDir(filePath) {
		return fileStrAsSlice[0], filePath
	}
	return fileStrAsSlice[0], strings.Join(filePathAsSlice[:len(filePathAsSlice)-1], "/")
}

func getASMFileName(path, fileName string) string {
	return path + "/" + fileName + ASM_EXTENSION
}

func getOutputFileFromInputPath(filePath string) (*os.File, string, error) {
	fileName, path := getFileNameAndDirectory(filePath)
	fileOutName := getASMFileName(path, fileName)
	fileOut, err := os.Create(fileOutName)
	if err != nil {
		return nil, "", errors.New(fmt.Sprintf("Could not make output file: %v", fileOutName))
	}
	return fileOut, fileName, nil
}

func getCodeWriter(filePath string) *CodeWriter {
	outputfile, fileName, err := getOutputFileFromInputPath(filePath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &CodeWriter{outputFile: outputfile, num: 0, fileName: fileName}
}

func getInit(num int) string {
	init := []string{}
	// SP=256
	init = append(init, "@256")
	init = append(init, "D=A")
	init = append(init, "@SP")
	init = append(init, "M=D")
	// Call Sys.init
	functionName := "Sys.init"
	returnAddress := getCallReturnAddress(functionName, num)
	init = append(init, strings.TrimSuffix(getCall(returnAddress, 0, functionName), "\n"))
	return joinStrings(init)
}

func (cw *CodeWriter) writeInit() {
	cw.writeStringToOutput(getInit(cw.incrementNum()))
}

func (cw *CodeWriter) setFileName(filePath string) {
	cw.num = 0
	fileName, _ := getFileNameAndDirectory(filePath)
	cw.fileName = fileName
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

func setComparison(cwNum int, op []string, comp string) []string {
	op = append(op, "D=M-D")
	// Jump to label if conditional is true
	op = append(op, fmt.Sprintf("@j%d", cwNum))
	op = append(op, fmt.Sprintf("D;%s", comp))
	// false
	op = append(op, "@SP")
	op = append(op, "A=M")
	op = append(op, fmt.Sprintf("M=%d", FALSE))
	op = append(op, fmt.Sprintf("@j%dend", cwNum))
	op = append(op, "0;JMP")
	// true
	op = append(op, fmt.Sprintf("(j%d)", cwNum))
	op = append(op, "@SP")
	op = append(op, "A=M")
	op = append(op, fmt.Sprintf("M=%d", TRUE))
	op = append(op, fmt.Sprintf("(j%dend)", cwNum))
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

func getEq(cwNum int) string {
	eq := []string{}
	eq = decrementStackPointer(eq)
	// D=*SP
	eq = append(eq, "D=M")
	eq = decrementStackPointer(eq)
	eq = setComparison(cwNum, eq, "JEQ")
	eq = incrementStackPointer(eq)
	return joinStrings(eq)
}

func getGt(cwNum int) string {
	gt := []string{}
	gt = decrementStackPointer(gt)
	// D=*SP
	gt = append(gt, "D=M")
	gt = decrementStackPointer(gt)
	gt = setComparison(cwNum, gt, "JGT")
	gt = incrementStackPointer(gt)
	return joinStrings(gt)
}

func getLt(cwNum int) string {
	lt := []string{}
	lt = decrementStackPointer(lt)
	// D=*SP
	lt = append(lt, "D=M")
	lt = decrementStackPointer(lt)
	lt = setComparison(cwNum, lt, "JLT")
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

func (cw *CodeWriter) writeArithmetic(command string) {
	switch command {
	case "add":
		cw.writeStringToOutput(getAdd())
	case "sub":
		cw.writeStringToOutput(getSub())
	case "neg":
		cw.writeStringToOutput(getNeg())
	case "eq":
		cw.writeStringToOutput(getEq(cw.incrementNum()))
	case "gt":
		cw.writeStringToOutput(getGt(cw.incrementNum()))
	case "lt":
		cw.writeStringToOutput(getLt(cw.incrementNum()))
	case "and":
		cw.writeStringToOutput(getAnd())
	case "or":
		cw.writeStringToOutput(getOr())
	case "not":
		cw.writeStringToOutput(getNot())
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

func getLocalPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(LOCAL_ABBR, idx)
	}
	return segmentPointerPop(LOCAL_ABBR, idx)
}

func getArgumentPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(ARGUMENT_ABBR, idx)
	}
	return segmentPointerPop(ARGUMENT_ABBR, idx)
}

func getThisPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(THIS_ABBR, idx)
	}
	return segmentPointerPop(THIS_ABBR, idx)
}

func getThatPushPop(pushPop int, idx int) string {
	if pushPop == C_PUSH {
		return segmentPointersPush(THAT_ABBR, idx)
	}
	return segmentPointerPop(THAT_ABBR, idx)
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

func getStaticPushPop(cw *CodeWriter, pushPop int, idx int) string {
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
	push = append(push, "A=D+A")
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

func getTempPushPop(pushPop int, idx int) string {
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

func getPointerPushPop(pushPop int, idx int) string {
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
		cw.writeStringToOutput(getLocalPushPop(pushPop, idx))
	case "argument":
		cw.writeStringToOutput(getArgumentPushPop(pushPop, idx))
	case "this":
		cw.writeStringToOutput(getThisPushPop(pushPop, idx))
	case "that":
		cw.writeStringToOutput(getThatPushPop(pushPop, idx))
	case "constant":
		cw.writeStringToOutput(getConstantPush(idx))
	case "static":
		cw.writeStringToOutput(getStaticPushPop(cw, pushPop, idx))
	case "temp":
		cw.writeStringToOutput(getTempPushPop(pushPop, idx))
	case "pointer":
		cw.writeStringToOutput(getPointerPushPop(pushPop, idx))
	}
}

func getCallReturnAddress(functionName string, num int) string {
	return fmt.Sprintf("%s$ret.%d", functionName, num)
}

func pushSegmentPointerAddress(push []string, abbr string) []string {
	push = append(push, fmt.Sprintf("@%s", abbr))
	push = append(push, "D=M")
	push = append(push, "@SP")
	push = append(push, "A=M")
	push = append(push, "M=D")
	push = incrementStackPointer(push)
	return push
}

func getLabel(label string) string {
	return fmt.Sprintf("(%s)\n", label)
}

func (cw *CodeWriter) writeLabel(label string) {
	cw.writeStringToOutput(getLabel(label))
}

func getGoto(label string) string {
	gt := []string{}
	gt = append(gt, fmt.Sprintf("@%s", label))
	gt = append(gt, fmt.Sprintf("0;JMP"))
	return joinStrings(gt)
}

func (cw *CodeWriter) writeGoto(label string) {
	cw.writeStringToOutput(getLabel(label))
}

func getIfGoto(label string, num int) string {
	ifgt := []string{}
	// D=*(SP-1)
	ifgt = decrementStackPointer(ifgt)
	ifgt = append(ifgt, "D=M")
	// if D == false jump to end
	ifLabel := fmt.Sprintf("ifgoto.%d", num)
	ifgt = append(ifgt, fmt.Sprintf("@%s", ifLabel))
	ifgt = append(ifgt, "D; JEQ")
	// if D == true jump to label
	ifgt = append(ifgt, fmt.Sprintf("@%s", label))
	ifgt = append(ifgt, "0; JMP")
	ifgt = append(ifgt, fmt.Sprintf("(%s)", ifLabel))
	return joinStrings(ifgt)
}

func (cw *CodeWriter) writeIf(label string) {
	cw.writeStringToOutput(getIfGoto(label, cw.incrementNum()))
}

func getCall(returnAddress string, nArgs int, functionName string) string {
	// push returnAddress
	call := []string{}
	call = append(call, fmt.Sprintf("@%s", returnAddress))
	call = append(call, "D=A")
	call = append(call, "@SP")
	call = append(call, "A=M")
	call = append(call, "M=D")
	call = incrementStackPointer(call)
	// push LCL
	call = pushSegmentPointerAddress(call, LOCAL_ABBR)
	// push ARG
	call = pushSegmentPointerAddress(call, ARGUMENT_ABBR)
	// push THIS
	call = pushSegmentPointerAddress(call, THIS_ABBR)
	// push THAT
	call = pushSegmentPointerAddress(call, THAT_ABBR)
	// ARG = SP - 5 - nArgs
	call = append(call, "@SP")
	call = append(call, "AD=M")
	call = append(call, "M=D")
	call = incrementStackPointer(call)
	call = append(call, strings.TrimSuffix(getConstantPush(5), "\n"))
	call = append(call, strings.TrimSuffix(getConstantPush(nArgs), "\n"))
	call = append(call, strings.TrimSuffix(getSub(), "\n"))
	call = append(call, strings.TrimSuffix(getSub(), "\n"))
	call = decrementStackPointer(call)
	call = append(call, "D=M")
	call = append(call, "@ARG")
	call = append(call, "M=D")
	// LCL = SP
	call = append(call, "@SP")
	call = append(call, "D=M")
	call = append(call, "@LCL")
	call = append(call, "M=D")
	// goto functionName
	call = append(call, fmt.Sprintf("@%s", functionName))
	call = append(call, "0;JMP")
	// (returnAddress)
	call = append(call, fmt.Sprintf("(%s)", returnAddress))
	return joinStrings(call)
}

func (cw *CodeWriter) writeCall(functionName string, numArgs int) {
	returnAddress := getCallReturnAddress(functionName, cw.incrementNum())
	cw.writeStringToOutput(getCall(returnAddress, numArgs, functionName))
}

func getFunction(functionName string, nVars int) string {
	function := []string{}
	// (functionName)
	function = append(function, fmt.Sprintf("(%s)", functionName))
	for i := 0; i < nVars; i++ {
		// push 0
		function = append(function, "@0")
		function = append(function, "D=A")
		function = append(function, "@SP")
		function = append(function, "A=M")
		function = append(function, "M=D")
		function = incrementStackPointer(function)
	}
	return joinStrings(function)
}

func (cw *CodeWriter) writeFunction(functionName string, nVars int) {
	cw.writeStringToOutput(getFunction(functionName, nVars))
}

func getFromEndFrame(instr []string, pointer string, offset int) []string {
	instr = append(instr, "@R13")
	instr = append(instr, "D=M")
	instr = append(instr, "@SP")
	instr = append(instr, "A=M")
	instr = append(instr, "M=D")
	instr = incrementStackPointer(instr)
	instr = append(instr, strings.TrimSuffix(getConstantPush(offset), "\n"))
	instr = append(instr, strings.TrimSuffix(getSub(), "\n"))
	instr = decrementStackPointer(instr)
	instr = append(instr, "D=M")
	instr = append(instr, "@D")
	instr = append(instr, "D=M")
	instr = append(instr, fmt.Sprintf("@%s", pointer))
	instr = append(instr, "M=D")
	return instr
}

func getReturn() string {
	rtrn := []string{}
	// endFrame = LCL
	rtrn = append(rtrn, "@LCL")
	rtrn = append(rtrn, "D=M")
	rtrn = append(rtrn, "@R13")
	rtrn = append(rtrn, "M=D")
	// retAddr = *(endFrame - 5)
	rtrn = getFromEndFrame(rtrn, "R14", 5)
	// *ARG = POP()
	rtrn = decrementStackPointer(rtrn)
	rtrn = append(rtrn, "D=M")
	rtrn = append(rtrn, "@ARG")
	rtrn = append(rtrn, "A=M")
	rtrn = append(rtrn, "M=D")
	// SP = ARG + 1
	rtrn = append(rtrn, "@ARG")
	rtrn = append(rtrn, "D=M")
	rtrn = append(rtrn, "@SP")
	rtrn = append(rtrn, "M=D+1")
	// THAT = *(endFrame - 1)
	rtrn = getFromEndFrame(rtrn, "THAT", 1)
	// THIS = *(endFrame - 2)
	rtrn = getFromEndFrame(rtrn, "THIS", 2)
	// ARG = *(endFrame - 3)
	rtrn = getFromEndFrame(rtrn, "ARG", 3)
	// LCL = *(endFrame - 4)
	rtrn = getFromEndFrame(rtrn, "LCL", 4)
	// goto retAddr
	rtrn = append(rtrn, "@R14")
	rtrn = append(rtrn, "0;JMP")
	return joinStrings(rtrn)
}

func (cw *CodeWriter) writeReturn() {
	cw.writeStringToOutput(getReturn())
}

func closeLoop() string {
	close := []string{}
	close = append(close, "(END_EXECUTION)")
	close = append(close, "@END_EXECUTION")
	close = append(close, "0;JMP")
	return joinStrings(close)
}

func (cw *CodeWriter) Close() {
	cw.writeStringToOutput(closeLoop())
	cw.outputFile.Close()
}
