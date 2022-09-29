package main

import (
	"fmt"
	"log"
	"strconv"
)

// AType => @[symbol]
// CType => dest=comp;jump

var COMP = map[string]string{
	"0":   "101010",
	"1":   "111111",
	"-1":  "111010",
	"D":   "001100",
	"A":   "110000",
	"M":   "110000",
	"!D":  "001101",
	"!A":  "110001",
	"!M":  "110001",
	"-D":  "001111",
	"-A":  "110011",
	"-M":  "110011",
	"D+1": "011111",
	"A+1": "110111",
	"M+1": "110111",
	"D-1": "001110",
	"A-1": "110010",
	"M-1": "110010",
	"D+A": "000010",
	"D+M": "000010",
	"D-A": "010011",
	"D-M": "010011",
	"A-D": "000111",
	"M-D": "000111",
	"D&A": "000000",
	"D&M": "000000",
	"D|A": "010101",
	"D|M": "010101",
}

var COMP_A = map[string]bool{
	"M":   true,
	"!M":  true,
	"-M":  true,
	"M+1": true,
	"M-1": true,
	"D+M": true,
	"D-M": true,
	"M-D": true,
	"D&M": true,
	"D|M": true,
}

var DEST = map[string]string{
	"null": "000",
	"M":    "001",
	"D":    "010",
	"MD":   "011",
	"A":    "100",
	"AM":   "101",
	"AD":   "110",
	"AMD":  "111",
}

var JUMP = map[string]string{
	"null": "000",
	"JGT":  "001",
	"JEQ":  "010",
	"JGE":  "011",
	"JLT":  "100",
	"JNE":  "101",
	"JLE":  "110",
	"JMP":  "111",
}

func compAIsSet(comp string) bool {
	_, found := COMP_A[comp]
	return found
}

func setA(comp string) string {
	if compAIsSet(comp) {
		return "1"
	}
	return "0"
}

func getATypeBinary(field []string, st *symbolTable) string {
	if val, found := st.getValue(field[0]); found {
		return fmt.Sprintf("%.16b\n", val)
	}
	digit, err := strconv.Atoi(field[0])
	if err != nil {
		log.Fatal("Error: Could not convert non digit string")
	}
	return fmt.Sprintf("%.16b\n", digit)
}

func getCTypeBinary(field []string, st *symbolTable) string {
	fieldLength := len(field)
	if fieldLength == 3 {
		// dest=comp;jump
		return fmt.Sprintf("111%v%v%v%v\n", setA(field[1]), COMP[field[1]], DEST[field[0]], JUMP[field[2]])
	} else {
		_, found := JUMP[field[1]]
		if found {
			// comp;jump
			return fmt.Sprintf("111%v%v%v%v\n", setA(field[0]), COMP[field[0]], DEST["null"], JUMP[field[1]])
		} else {
			// dest=comp
			return fmt.Sprintf("111%v%v%v%v\n", setA(field[1]), COMP[field[1]], DEST[field[0]], JUMP["null"])
		}
	}
}

func getBinary(ft *fieldTuple, st *symbolTable) string {
	if ft.cmd == AType {
		return getATypeBinary(ft.fields, st)
	}
	return getCTypeBinary(ft.fields, st)
}
