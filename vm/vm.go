package vm

import "fmt"

var pc int64

type VM struct {
	pc        int64   // Program counter
	bc        []int64 // Bytecode (program)
	registers []int64 // General purpose registers
	stack     *stack  // Program stack
}

func New(in []int64) *VM {
	return &VM{
		bc:        in,
		registers: make([]int64, totalRegisters),
		stack:     &stack{},
	}
}

func (vm *VM) Start() int64 {
	for {
		code := vm.fetch()

		switch code {
		case Halt:
			return vm.fetch()
		case Push:
			vm.stack.push(vm.fetch())
		case PushReg:
			vm.stack.push(vm.registers[vm.fetch()])
		case Pop:
			vm.stack.pop()
		case PopReg:
			vm.registers[vm.fetch()] = vm.stack.pop()
		case Add:
			right := vm.stack.pop()
			left := vm.stack.pop()
			vm.stack.push(left + right)
		case Sub:
			right := vm.stack.pop()
			left := vm.stack.pop()
			vm.stack.push(left - right)
		case Set:
			reg := vm.fetch()
			vm.registers[reg] = vm.fetch()
		case Jump:
			vm.pc = vm.fetch()
		case JumpGtz:
			tos := vm.stack.tos()
			next := vm.fetch()
			if tos > 0 {
				vm.pc = next
			}
		case Print:
			fmt.Println(vm.stack.tos())
		}
	}
}

func (vm *VM) fetch() int64 {
	c := vm.bc[vm.pc]
	vm.pc++
	return c
}
