package lexer

import "github.com/elemental-vm/test-vm/vm"

var bytecodes = map[string]int64{
	"HALT":    vm.Halt,
	"EXIT":    vm.Halt,
	"PUSH":    vm.Push,
	"PUSHREG": vm.PushReg,
	"POP":     vm.Pop,
	"POPREG":  vm.PopReg,
	"ADD":     vm.Add,
	"SUB":     vm.Sub,
	"SET":     vm.Set,
	"JMP":     vm.Jump,
	"JMPGZ":   vm.JumpGtz,
	"PRINT":   vm.Print,
}

var registers = map[string]int64{
	"A": vm.A,
	"B": vm.B,
	"C": vm.C,
	"D": vm.D,
	"E": vm.E,
	"F": vm.F,
	"G": vm.G,
	"H": vm.H,
	"I": vm.I,
	"J": vm.J,
}
