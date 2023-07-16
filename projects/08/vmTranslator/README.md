# Running
Run the vm translator using command: 

`go run main.go parser.go codeWriter.go [path to .asm file]` 

# Notes
Some of the tests may not pass due to the codeWriter.writeInit() function
which is inconsistently tested which will result in incorrect test
results with given tests. If a tests does not pass, comment out this line
in main.go and rerun with outputted .asm.