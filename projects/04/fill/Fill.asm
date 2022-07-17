// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Save screen addr
@SCREEN
D = A
@screenAddr
M = D
// Set screen dimensions
@32
D = A
@screenWidth
M = D
@256
D = A
@screenHeight
M = D
// Set colors
@white
M = 0
@black
M = -1
// Save kbd addr
@KBD
D = A
@kbdAddr
M = D
// Constantly check if key is pressed
(LOOP)
// Pixel value initially set to clear screen
@white
D = M
@pixelValue
M = D
@kbdAddr
A = M
D = M
@POPULATE_SCREEN
D; JEQ
// Darken screen if key is pressed
@black
D = M
@pixelValue
M = D
(POPULATE_SCREEN)
// Check if pixel value is already set to desired state
// don't populate if it is
@pixelValue
D = M
@screenAddr
A = M
D = M
@pixelValue
D = D - M
@FINISHED_POPULATING
D; JEQ
// Fill screen
@screenAddr
D = M
@currentScreenAddr
M = D
@0
D = A
@i
M = D
(POPULATE_COLUMN)
  @i
  D = M
  @screenWidth
  D = D - M
  @FINISHED_POPULATING
  D; JEQ
  @0
  D = A
  @j
  M = D
  (POPULATE_ROW)
    @j
    D = M
    @screenHeight
    D = D - M
    @FINISHED_FILLING_ROW
    D; JEQ
    @pixelValue
    D = M
    @currentScreenAddr
    A = M
    M = D
    @currentScreenAddr
    M = M + 1
    @j
    M = M + 1
    @POPULATE_ROW
    0; JMP
  (FINISHED_FILLING_ROW)
  @i
  M = M + 1
@POPULATE_COLUMN
0; JMP
(FINISHED_POPULATING)
@LOOP
0; JMP