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
		//fmt.Printf("Executing %d\n", code)

		switch code {
		case Halt:
			return vm.fetch()
		case Push:
			vm.stack.push(vm.fetch())
		case PushReg:
			vm.stack.push(vm.registers[vm.fetch()])
		case Dup:
			vm.stack.push(vm.stack.tos())
		case Pop:
			vm.stack.pop()
		case PopReg:
			vm.registers[vm.fetch()] = vm.stack.pop()
		case Swap:
			a := vm.stack.pop()
			b := vm.stack.pop()
			vm.stack.push(a)
			vm.stack.push(b)
		case Add:
			right := vm.stack.pop()
			left := vm.stack.pop()
			vm.stack.push(left + right)
		case Sub:
			right := vm.stack.pop()
			left := vm.stack.pop()
			vm.stack.push(left - right)
		case Mul:
			right := vm.stack.pop()
			left := vm.stack.pop()
			vm.stack.push(left * right)
		case Div:
			right := vm.stack.pop()
			left := vm.stack.pop()
			vm.stack.push(left / right)
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
		case JumpLtz:
			tos := vm.stack.tos()
			next := vm.fetch()
			if tos < 0 {
				vm.pc = next
			}
		case JumpEq:
			tos := vm.stack.tos()
			next := vm.fetch()
			if tos == 0 {
				vm.pc = next
			}
		case JumpNeq:
			tos := vm.stack.tos()
			next := vm.fetch()
			if tos != 0 {
				vm.pc = next
			}
		case Print:
			fmt.Println(vm.stack.tos())
		case PrintS:
			fmt.Println(vm.stack.string())
		case Return:
			vm.pc = vm.stack.pop()
		case Call:
			fn := vm.fetch()
			vm.stack.push(vm.pc)
			vm.pc = fn
		}
	}
}

func (vm *VM) fetch() int64 {
	c := vm.bc[vm.pc]
	vm.pc++
	return c
}
