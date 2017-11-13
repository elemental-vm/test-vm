package lexer

import "github.com/elemental-vm/test-vm/vm"

var bytecodes = map[string]byte{
	"HALT": vm.Halt,
	"EXIT": vm.Halt,

	"PUSHI":   vm.PushI,
	"PUSHSTR": vm.PushStr,
	"PUSHREG": vm.PushReg,
	"SWAP":    vm.Swap,
	"DUP":     vm.Dup,

	"POP":    vm.Pop,
	"POPREG": vm.PopReg,
	"STORE":  vm.Store,

	"ADD": vm.Add,
	"SUB": vm.Sub,
	"MUL": vm.Mul,
	"DIV": vm.Div,

	"SETI":   vm.SetI,
	"SETSTR": vm.SetStr,

	"JMP":    vm.Jump,
	"JMPGZ":  vm.JumpGtz,
	"JMPLZ":  vm.JumpLtz,
	"JMPEQ":  vm.JumpEq,
	"JMPNEQ": vm.JumpNeq,
	"JMPREG": vm.JumpReg,

	"PRINT":  vm.Print,
	"DUMP":   vm.Dump,
	"PRINTR": vm.PrintR,
	"DUMPR":  vm.DumpR,

	"RETURN": vm.Return,
	"CALL":   vm.Call,

	"CONCAT": vm.Concat,
	"PARAM":  vm.Param,

	"CMP": vm.Compare,

	"JMPZGZ":  vm.JumpZGtz,
	"JMPZLZ":  vm.JumpZLtz,
	"JMPZEQ":  vm.JumpZEq,
	"JMPZNEQ": vm.JumpZNeq,

	"STEP": vm.Step,
}

var registers = map[string]byte{
	"RT": vm.RT,
	"SP": vm.SP,
	"FP": vm.FP,
	"PC": vm.PC,
	"A":  vm.A,
	"B":  vm.B,
	"C":  vm.C,
	"D":  vm.D,
	"E":  vm.E,
	"F":  vm.F,
	"G":  vm.G,
	"H":  vm.H,
	"I":  vm.I,
	"J":  vm.J,
}
