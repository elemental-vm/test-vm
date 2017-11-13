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

	// PC is the current program counter register
	PC = totalUserRegisters
	// SP is the current stack pointer register
	SP = totalUserRegisters + 1 // (sp - 1 == TOS value)
	// FP is the current frame pointer register
	FP = totalUserRegisters + 2
	// RT is the current return address register
	RT = totalUserRegisters + 3

	totalRegisters = RT + 1
)

type vmValue struct {
	t    regType
	iVal int64
	sVal []byte
}

func (v *vmValue) dup() *vmValue {
	return &vmValue{
		t:    v.t,
		iVal: v.iVal,
		sVal: v.sVal,
	}
}

type VM struct {
	flags struct {
		debug bool
		zero  int8
	}
	errorMsg string

	program []byte // Bytecode (program)

	registers []*vmValue // General purpose registers
	stack     []*vmValue // Stack
}

func New(in []byte) *VM {
	vm := &VM{
		program:   in,
		registers: make([]*vmValue, totalRegisters),
		stack:     make([]*vmValue, 1024),
	}

	for i := range vm.registers {
		vm.registers[i] = &vmValue{}
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
		case PopReg:
			vm.opPopReg()
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
		case JumpZGtz:
			vm.opJumpZGtz()
		case JumpZLtz:
			vm.opJumpZLtz()
		case JumpZEq:
			vm.opJumpZEq()
		case JumpZNeq:
			vm.opJumpZNeq()
		case JumpReg:
			vm.opJumpReg()

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

		case Concat:
			vm.opConcat()

		case Param:
			vm.opParam()

		case Compare:
			vm.opCompare()

		default:
			fmt.Printf("Unknown bytecode 0x%X\n", code)
			return 1
		}
	}
}

func (vm *VM) fetch() byte {
	nextpc := vm.registers[PC].iVal
	if nextpc >= int64(len(vm.program)) {
		panic("Reached end of program with no HALT instruction, halting")
	}
	c := vm.program[nextpc]
	vm.registers[PC].iVal++
	return c
}

func (vm *VM) setPC(v int64) {
	vm.registers[PC].iVal = v
}

func (vm *VM) getPC() int64 {
	return vm.registers[PC].iVal
}

func (vm *VM) pushStack(v *vmValue) {
	vm.stack[vm.registers[SP].iVal] = v
	vm.registers[SP].iVal++
}

func (vm *VM) pushStackI(v int64) {
	csp := vm.registers[SP].iVal
	if vm.stack[csp] == nil {
		vm.stack[csp] = &vmValue{}
	}

	vm.stack[csp].t = regInt
	vm.stack[csp].iVal = v
	vm.registers[SP].iVal++
}

func (vm *VM) pushStackStr(v []byte) {
	csp := vm.registers[SP].iVal
	if vm.stack[csp] == nil {
		vm.stack[csp] = &vmValue{}
	}

	vm.stack[csp].t = regStr
	vm.stack[csp].sVal = v
	vm.registers[SP].iVal++
}

func (vm *VM) popStack() *vmValue {
	vm.registers[SP].iVal--
	sv := vm.stack[vm.registers[SP].iVal]
	return &vmValue{
		t:    sv.t,
		iVal: sv.iVal,
		sVal: sv.sVal,
	}
}

func (vm *VM) getTOS() *vmValue {
	sv := vm.stack[vm.registers[SP].iVal-1]
	return &vmValue{
		t:    sv.t,
		iVal: sv.iVal,
		sVal: sv.sVal,
	}
}

func (vm *VM) printStack() {
	if vm.registers[SP].iVal == 0 {
		fmt.Println("[]")
		return
	}

	sp := vm.registers[SP].iVal - 1
	var out bytes.Buffer
	out.WriteByte('[')

	for sp >= 0 {
		if vm.stack[sp].t == regInt {
			out.WriteString("0x")
			out.WriteString(strconv.FormatInt(vm.stack[sp].iVal, 16))
		} else {
			out.Write(vm.stack[sp].sVal)
		}
		if sp > 0 {
			out.WriteString(", ")
		}
		sp--
	}

	out.WriteByte(']')
	fmt.Println(out.String())
}

func (vm *VM) printRegisters() {
	i := byte(0)
	fmt.Printf("| PC: 0x%X | SP: 0x%X | FP: 0x%X | RT: 0x%X | ",
		vm.registers[PC].iVal,
		vm.registers[SP].iVal,
		vm.registers[FP].iVal,
		vm.registers[RT].iVal,
	)

	for i < totalUserRegisters {
		if vm.registers[i].t == regInt {
			fmt.Printf("%c: 0x%X | ", 'A'+i, vm.registers[i].iVal)
		} else {
			fmt.Printf("%c: %q | ", 'A'+i, vm.registers[i].sVal)
		}
		i++
	}
	fmt.Println("")
}
