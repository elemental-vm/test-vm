package vm

// Instruction set
const (
	Halt    = iota // 0
	Push           // 1
	PushReg        // 2
	Pop            // 3
	PopReg         // 4
	Add            // 5
	Sub            // 6
	Set            // 7
	Jump           // 8
	JumpGtz        // 9
	JumpLtz        // 10
	JumpEq         // 11
	JumpNeq        // 12
	Print          // 13
	Return         // 14
	Call           // 15
	PrintS         // 16
	Swap           // 17
	Dup            // 18
	Mul            // 19
	Div            // 20
)

// Registers
const (
	RT = iota // 0
	A         // 1
	B         // 2
	C         // 3
	D         // 4
	E         // 5
	F         // 6
	G         // 7
	H         // 8
	I         // 9
	J         // 10
	totalRegisters
)
