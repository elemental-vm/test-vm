package vm

// Instruction set
const (
	Halt byte = iota // 0x00

	PushI   // 0x01
	PushStr // 0x02
	PushReg // 0x03
	Pop     // 0x04
	PopReg  // 0x05
	Store   // 0x06
	Swap    // 0x07
	Dup     // 0x08

	Add // 0x09
	Sub // 0x0A
	Mul // 0x0B
	Div // 0x0C

	SetI   // 0x0D
	SetStr // 0x0E

	Jump    // 0x0F
	JumpGtz // 0x10
	JumpLtz // 0x11
	JumpEq  // 0x12
	JumpNeq // 0x13

	Print  // 0x14
	PrintR // 0x15
	Dump   // 0x16
	DumpR  // 0x17

	Return // 0x18
	Call   // 0x19

	Concat // 0x1A
)

// Registers
const (
	A              byte = iota // 0x00
	B                          // 0x01
	C                          // 0x02
	D                          // 0x03
	E                          // 0x04
	F                          // 0x05
	G                          // 0x06
	H                          // 0x07
	I                          // 0x08
	J                          // 0x09
	totalRegisters             // 10
	RT                         // 0x11
)
