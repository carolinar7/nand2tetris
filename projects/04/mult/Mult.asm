// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

@0
D = A
@R2
M = D
@R0
D = M
@FIRSTJUMP
D; JGE
@0
D = A - D
(FIRSTJUMP)
@R3
M = D
(THIRDJUMP)
@R3
D = M
@SECONDJUMP
D; JEQ
@R1
D = M
@R2
M = M + D
@R3
M = M - 1
@THIRDJUMP
0; JMP
(SECONDJUMP)
@R0
D = M
@END
D; JGT
@0
D = A
@R2
M = D - M
(END)
@END
0; JMP