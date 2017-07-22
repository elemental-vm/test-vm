package vm

import (
	"encoding/binary"
)

func (vm *VM) opPushI() {
	vm.pushStackI(vm.getInt64())
}
func (vm *VM) opPushStr() {
	vm.pushStackStr(vm.fetchString())
}
func (vm *VM) opPushReg() {
	reg := vm.registers[vm.fetch()]
	if reg.t == regInt {
		vm.pushStackI(reg.iVal)
	} else {
		vm.pushStackStr(reg.sVal)
	}
}
func (vm *VM) opDup() {
	vm.pushStack(vm.getTOS())
}
func (vm *VM) opPop() {
	vm.popStack()
}
func (vm *VM) opStore() {
	vm.registers[vm.fetch()] = vm.popStack()
}
func (vm *VM) opSwap() {
	vm.stack[vm.sp-2], vm.stack[vm.sp-1] = vm.stack[vm.sp-1], vm.stack[vm.sp-2]
}

func (vm *VM) opAdd() {
	right := vm.popStack()
	if right.t != regInt {
		vm.errorMsg = "ADD only works on integers"
		return
	}

	left := vm.popStack()
	if right.t != regInt {
		vm.errorMsg = "ADD only works on integers"
		return
	}

	vm.pushStackI(left.iVal + right.iVal)
}
func (vm *VM) opSub() {
	right := vm.popStack()
	if right.t != regInt {
		vm.errorMsg = "SUB only works on integers"
		return
	}

	left := vm.popStack()
	if right.t != regInt {
		vm.errorMsg = "SUB only works on integers"
		return
	}

	vm.pushStackI(left.iVal - right.iVal)
}
func (vm *VM) opMul() {
	right := vm.popStack()
	if right.t != regInt {
		vm.errorMsg = "MUL only works on integers"
		return
	}

	left := vm.popStack()
	if right.t != regInt {
		vm.errorMsg = "MUL only works on integers"
		return
	}

	vm.pushStackI(left.iVal * right.iVal)
}
func (vm *VM) opDiv() {
	right := vm.popStack()
	if right.t != regInt {
		vm.errorMsg = "DIV only works on integers"
		return
	}

	left := vm.popStack()
	if right.t != regInt {
		vm.errorMsg = "DIV only works on integers"
		return
	}

	vm.pushStackI(left.iVal / right.iVal)
}

func (vm *VM) opSetI() {
	reg := vm.fetch()
	vm.registers[reg].t = regInt
	vm.registers[reg].iVal = vm.getInt64()
}

func (vm *VM) opSetStr() {
	reg := vm.fetch()
	vm.registers[reg].t = regStr
	vm.registers[reg].sVal = vm.fetchString()
}

func (vm *VM) opJump() {
	vm.pc = vm.getInt64()
}
func (vm *VM) opJumpGtz() {
	next := vm.getInt64()
	if vm.getTOS().iVal > 0 {
		vm.pc = next
	}
}
func (vm *VM) opJumpLtz() {
	next := vm.getInt64()
	if vm.getTOS().iVal < 0 {
		vm.pc = next
	}
}
func (vm *VM) opJumpEq() {
	next := vm.getInt64()
	if vm.getTOS().iVal == 0 {
		vm.pc = next
	}
}
func (vm *VM) opJumpNeq() {
	next := vm.getInt64()
	if vm.getTOS().iVal != 0 {
		vm.pc = next
	}
}

func (vm *VM) opReturn() {
	vm.pc = vm.popStack().iVal
}
func (vm *VM) opCall() {
	fn := vm.getInt64()
	vm.pushStackI(vm.pc)
	vm.pc = fn
}

func (vm *VM) getInt64() int64 {
	buf := make([]byte, 8)
	buf[0] = vm.fetch()
	buf[1] = vm.fetch()
	buf[2] = vm.fetch()
	buf[3] = vm.fetch()
	buf[4] = vm.fetch()
	buf[5] = vm.fetch()
	buf[6] = vm.fetch()
	buf[7] = vm.fetch()

	i, _ := binary.Varint(buf)
	return i
}

func (vm *VM) fetchString() string {
	l := (int16(vm.fetch()) << 8) + int16(vm.fetch())
	str := make([]byte, l)

	for i := int16(0); i < l; i++ {
		str[i] = vm.fetch()
	}

	return string(str)
}
