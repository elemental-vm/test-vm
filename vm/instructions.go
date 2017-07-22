package vm

// Instruction set
const (
	Halt byte = iota // 0x00

	PushI   // 0x01
	PushStr // 0x02
	PushReg // 0x03
	Pop     // 0x04
	Store   // 0x05
	Swap    // 0x06
	Dup     // 0x07

	Add // 0x08
	Sub // 0x09
	Mul // 0x0A
	Div // 0x0B

	SetI   // 0x0C
	SetStr // 0x0D

	Jump    // 0x0E
	JumpGtz // 0x0F
	JumpLtz // 0x10
	JumpEq  // 0x11
	JumpNeq // 0x12

	Print  // 0x13
	PrintR // 0x14
	Dump   // 0x15
	DumpR  // 0x16

	Return // 0x17
	Call   // 0x18
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
