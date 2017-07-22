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
	Print          // 10
)

// Registers
const (
	A = iota // 0
	B        // 1
	C        // 2
	D        // 3
	E        // 4
	F        // 5
	G        // 6
	H        // 7
	I        // 8
	J        // 9
	totalRegisters
)
