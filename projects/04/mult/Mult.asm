// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

// Reset value in R2
@0
D = A
@R2
M = D
// Store |R0| as a counter in R3
@R0
D = M
@IS_NOT_NEGATIVE
D; JGE
@0
D = A - D
(IS_NOT_NEGATIVE)
@counter
M = D
// Multiplication logic: Add R1 value until R3 counter goes to zero to R2
(LOOP)
@counter
D = M
@COUNTER_IS_ZERO
D; JEQ
@R1
D = M
@R2
M = M + D
@counter
M = M - 1
@LOOP
0; JMP
(COUNTER_IS_ZERO)
// Check if R0 is negative and flip sign of R2 if it is
@R0
D = M
@END
D; JGT
@0
D = A
@R2
M = D - M
// Infinite loop
(END)
@END
0; JMP