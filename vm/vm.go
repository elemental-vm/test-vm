package vm

import (
	"bytes"
	"fmt"
	"strconv"
)

type regType uint8

const (
	regInt regType = iota
	regStr
)

type register struct {
	t    regType
	iVal int64
	sVal []byte
}

type VM struct {
	flags struct {
		debug bool
	}
	errorMsg string

	pc int64 // Program counter
	sp int64 // Stack pointer (sp - 1 == TOS value)

	program []byte // Bytecode (program)

	registers []*register // General purpose registers
	stack     []*register // Stack
}

func New(in []byte) *VM {
	vm := &VM{
		program:   in,
		registers: make([]*register, totalRegisters),
		stack:     make([]*register, 1024),
	}

	for i := range vm.registers {
		vm.registers[i] = &register{}
	}

	return vm
}

func (vm *VM) Start(debug bool) byte {
	vm.flags.debug = debug

	for {
		if vm.errorMsg != "" {
			fmt.Println(vm.errorMsg)
			return 1
		}

		code := vm.fetch()

		if vm.flags.debug {
			fmt.Printf("Executing 0x%X\n", code)
		}

		switch code {
		case Halt:
			return vm.fetch()

		case PushI:
			vm.opPushI()
		case PushStr:
			vm.opPushStr()
		case PushReg:
			vm.opPushReg()
		case Dup:
			vm.opDup()
		case Pop:
			vm.opPop()
		case Store:
			vm.opStore()
		case Swap:
			vm.opSwap()

		case Add:
			vm.opAdd()
		case Sub:
			vm.opSub()
		case Mul:
			vm.opMul()
		case Div:
			vm.opDiv()

		case SetI:
			vm.opSetI()
		case SetStr:
			vm.opSetStr()

		case Jump:
			vm.opJump()
		case JumpGtz:
			vm.opJumpGtz()
		case JumpLtz:
			vm.opJumpLtz()
		case JumpEq:
			vm.opJumpEq()
		case JumpNeq:
			vm.opJumpNeq()

		case Print:
			tos := vm.getTOS()
			if tos.t == regInt {
				fmt.Printf("%d\n", tos.iVal)
			} else {
				fmt.Printf("%q\n", tos.sVal)
			}
		case Dump:
			vm.printStack()
		case PrintR:
			reg := vm.fetch()
			if vm.registers[reg].t == regInt {
				fmt.Printf("%d\n", vm.registers[reg].iVal)
			} else {
				fmt.Printf("%q\n", vm.registers[reg].sVal)
			}
		case DumpR:
			vm.printRegisters()

		case Return:
			vm.opReturn()
		case Call:
			vm.opCall()

		default:
			fmt.Printf("Unknown bytecode 0x%X", code)
			return 1
		}
	}
}

func (vm *VM) fetch() byte {
	c := vm.program[vm.pc]
	vm.pc++
	return c
}

func (vm *VM) pushStack(v *register) {
	vm.stack[vm.sp] = v
	vm.sp++
}

func (vm *VM) pushStackI(v int64) {
	if vm.stack[vm.sp] == nil {
		vm.stack[vm.sp] = &register{}
	}

	vm.stack[vm.sp].t = regInt
	vm.stack[vm.sp].iVal = v
	vm.sp++
}

func (vm *VM) pushStackStr(v []byte) {
	if vm.stack[vm.sp] == nil {
		vm.stack[vm.sp] = &register{}
	}

	vm.stack[vm.sp].t = regStr
	vm.stack[vm.sp].sVal = v
	vm.sp++
}

func (vm *VM) popStack() *register {
	vm.sp--
	sv := vm.stack[vm.sp]
	return &register{
		t:    sv.t,
		iVal: sv.iVal,
		sVal: sv.sVal,
	}
}

func (vm *VM) getTOS() *register {
	sv := vm.stack[vm.sp-1]
	return &register{
		t:    sv.t,
		iVal: sv.iVal,
		sVal: sv.sVal,
	}
}

func (vm *VM) printStack() {
	if vm.sp == 0 {
		fmt.Println("[]")
		return
	}

	sp := vm.sp - 1
	var out bytes.Buffer
	out.WriteByte('[')

	for sp >= 0 {
		if vm.registers[sp].t == regInt {
			out.WriteString(strconv.FormatInt(vm.stack[sp].iVal, 10))
		} else {
			out.Write(vm.stack[sp].sVal)
		}
		if sp > 0 {
			out.WriteByte(',')
		}
		sp--
	}

	out.WriteByte(']')
	fmt.Println(out.String())
}

func (vm *VM) printRegisters() {
	i := byte(0)
	for i < totalRegisters {
		if vm.registers[i].t == regInt {
			fmt.Printf("%c: %d\n", 'A'+i, vm.registers[i].iVal)
		} else {
			fmt.Printf("%c: %q\n", 'A'+i, vm.registers[i].sVal)
		}
		i++
	}
}
