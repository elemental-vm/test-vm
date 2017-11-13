package lexer

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"bytes"

	"io/ioutil"

	"github.com/elemental-vm/test-vm/vm"
)

var FileHeader = []byte{31, 'E', 'B', 'C'}

type sub struct {
	pos   int64
	label string
}

type Lexer struct {
	r      *bufio.Reader
	file   *os.File
	simple bool

	line int
	pc   int64

	program []byte

	labels    map[string]int64
	labelSubs []*sub
}

func New(file string) (*Lexer, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	l := &Lexer{}

	header := make([]byte, 4)
	f.Read(header)
	if bytes.Equal(header, FileHeader) {
		l.simple = true
		l.file = f
		return l, nil
	}

	f.Seek(0, 0) // Reset reader
	l.r = bufio.NewReader(f)
	l.labels = make(map[string]int64)
	l.labelSubs = make([]*sub, 0, 15)
	return l, nil
}

func (l *Lexer) addToProgram(bit byte) {
	l.pc++
	l.program = append(l.program, bit)
}

func (l *Lexer) addSliceToProgram(bytes []byte) {
	for _, b := range bytes {
		l.addToProgram(b)
	}
}

func (l *Lexer) addLabelSub(label string) {
	l.labelSubs = append(l.labelSubs, &sub{
		pos:   l.pc,
		label: label,
	})
}

func (l *Lexer) Parse() ([]byte, error) {
	if l.simple {
		program, err := ioutil.ReadAll(l.file)
		l.file.Close()
		if err != nil {
			return nil, err
		}
		return program, nil
	}

	l.program = make([]byte, 0, 1024)
	quit := false

	for {
		if quit {
			break
		}

		line, err := l.r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			quit = true
		}
		l.line++

		line = strings.TrimSpace(line)
		if line == "" || line[0] == ';' {
			continue
		}

		text := strings.SplitN(line, ";", 2)
		line = strings.TrimSpace(text[0])

		lineLabel := strings.Split(line, ":")
		if len(lineLabel) > 1 {
			l.labels[lineLabel[0]] = l.pc
			line = strings.TrimSpace(lineLabel[1])
		}

		structure := strings.Split(line, " ")
		parts := len(structure)
		if parts == 0 {
			continue
		}

		bytecode, ok := bytecodes[strings.ToUpper(structure[0])]
		if !ok {
			return nil, fmt.Errorf("Unknown instruction %s on line %d", structure[0], l.line)
		}
		l.addToProgram(bytecode)

		switch bytecode {
		case vm.Halt:
			err = l.parseParamOneByte(structure)
		case vm.PushI:
			err = l.parseParamOneInt(structure)
		case vm.PushReg:
			err = l.parseParamOneRegister(structure)
		case vm.PrintR:
			err = l.parseParamOneRegister(structure)
		case vm.PopReg:
			err = l.parseParamOneRegister(structure)
		case vm.Store:
			err = l.parseParamOneRegister(structure)
		case vm.SetI:
			err = l.parseParamsRegIntOrLabel(structure)
		case vm.Jump:
			err = l.parseParamOneIntOrLabel(structure)
		case vm.JumpGtz:
			err = l.parseParamOneIntOrLabel(structure)
		case vm.JumpLtz:
			err = l.parseParamOneIntOrLabel(structure)
		case vm.JumpEq:
			err = l.parseParamOneIntOrLabel(structure)
		case vm.JumpNeq:
			err = l.parseParamOneIntOrLabel(structure)
		case vm.Call:
			err = l.parseParamOneIntOrLabel(structure)
		case vm.PushStr:
			err = l.parseParamOneString(structure)
		case vm.SetStr:
			err = l.parseParamsRegString(structure)
		case vm.Param:
			err = l.parseParamsRegInt(structure)
		case vm.JumpReg:
			err = l.parseParamOneRegister(structure)
		case vm.Compare:
			err = l.parseParamsTwoRegisters(structure)
		case vm.JumpZGtz:
			err = l.parseParamOneIntOrLabel(structure)
		case vm.JumpZLtz:
			err = l.parseParamOneIntOrLabel(structure)
		case vm.JumpZEq:
			err = l.parseParamOneIntOrLabel(structure)
		case vm.JumpZNeq:
			err = l.parseParamOneIntOrLabel(structure)
		default:
			err = nil
		}

		if err != nil {
			return nil, err
		}
	}

	if err := l.subLabels(); err != nil {
		return nil, err
	}
	return l.program, nil
}

func (l *Lexer) parseParamOneByte(structure []string) error {
	if len(structure) != 2 {
		return fmt.Errorf("Expected int on line %d", l.line)
	}

	code, err := strconv.ParseInt(structure[1], 10, 64)
	if err != nil {
		return err
	}

	if code > 255 || code < 0 {
		return fmt.Errorf("Exit code must be between 0-255, line %d", l.line)
	}

	l.addToProgram(byte(code))
	return nil
}

func (l *Lexer) parseParamOneInt(structure []string) error {
	if len(structure) != 2 {
		return fmt.Errorf("Expected int on line %d", l.line)
	}

	code, err := strconv.ParseInt(structure[1], 10, 64)
	if err != nil {
		return err
	}
	l.addSliceToProgram(intToBytes(code))
	return nil
}

func (l *Lexer) parseParamOneIntOrLabel(structure []string) error {
	if len(structure) != 2 {
		return fmt.Errorf("Expected int on line %d", l.line)
	}

	if structure[1][0] == '%' {
		l.addLabelSub(structure[1][1:]) // Locations are 64 bits
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		return nil
	}

	code, err := strconv.ParseInt(structure[1], 10, 64)
	if err != nil {
		return err
	}
	l.addSliceToProgram(intToBytes(code))
	return nil
}

func (l *Lexer) parseParamsRegIntOrLabel(structure []string) error {
	if len(structure) != 3 {
		return fmt.Errorf("Expected int on line %d", l.line)
	}

	if structure[1][0] != '$' {
		return fmt.Errorf("Expected register on line %d", l.line)
	}

	reg, ok := getRegister(structure[1][1:])
	if !ok {
		return fmt.Errorf("%s is not a register", structure[1])
	}
	l.addToProgram(reg)

	if structure[2][0] == '%' {
		l.addLabelSub(structure[2][1:]) // Locations are 64 bits
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		l.addToProgram(0)
		return nil
	}

	code, err := strconv.ParseInt(structure[2], 10, 64)
	if err != nil {
		return err
	}
	l.addSliceToProgram(intToBytes(code))
	return nil
}

func (l *Lexer) parseParamOneRegister(structure []string) error {
	if len(structure) != 2 {
		return fmt.Errorf("Expected register on line %d", l.line)
	}

	if structure[1][0] != '$' {
		return fmt.Errorf("Expected register on line %d", l.line)
	}

	reg, ok := getRegister(structure[1][1:])
	if !ok {
		return fmt.Errorf("%s is not a register", structure[1])
	}
	l.addToProgram(reg)
	return nil
}

func (l *Lexer) parseParamsTwoRegisters(structure []string) error {
	if len(structure) != 3 {
		return fmt.Errorf("Expected register on line %d", l.line)
	}

	// First register
	if structure[1][0] != '$' {
		return fmt.Errorf("Expected register on line %d", l.line)
	}

	reg, ok := getRegister(structure[1][1:])
	if !ok {
		return fmt.Errorf("%s is not a register", structure[1])
	}
	l.addToProgram(reg)

	// Second register
	if structure[2][0] != '$' {
		return fmt.Errorf("Expected register on line %d", l.line)
	}

	reg, ok = getRegister(structure[2][1:])
	if !ok {
		return fmt.Errorf("%s is not a register", structure[2])
	}
	l.addToProgram(reg)
	return nil
}

func (l *Lexer) parseParamsRegInt(structure []string) error {
	if len(structure) != 3 {
		return fmt.Errorf("Expected register and int on line %d", l.line)
	}

	if structure[1][0] != '$' {
		return fmt.Errorf("Expected register on line %d", l.line)
	}

	// Register
	reg, ok := getRegister(structure[1][1:])
	if !ok {
		return fmt.Errorf("%s is not a register", structure[1])
	}
	l.addToProgram(reg)

	// Value
	code, err := strconv.ParseInt(structure[2], 10, 64)
	if err != nil {
		return err
	}
	l.addSliceToProgram(intToBytes(code))
	return nil
}

func (l *Lexer) parseParamOneString(structure []string) error {
	if len(structure) < 2 {
		return fmt.Errorf("Expected string on line %d", l.line)
	}

	str := strings.Join(structure[1:], " ")
	str = str[1:]
	str = str[:len(str)-1]

	strLen := len(str)
	if strLen > 32768 {
		return fmt.Errorf("String too long on line %d", l.line)
	}

	// Push string lenth
	l.addToProgram(byte(strLen >> 8))
	l.addToProgram(byte(strLen))

	// Add string literal
	l.addSliceToProgram([]byte(str))
	return nil
}

func (l *Lexer) parseParamsRegString(structure []string) error {
	if len(structure) < 3 {
		return fmt.Errorf("Expected string on line %d", l.line)
	}

	if structure[1][0] != '$' {
		return fmt.Errorf("Expected register on line %d", l.line)
	}

	// Register
	reg, ok := getRegister(structure[1][1:])
	if !ok {
		return fmt.Errorf("%s is not a register", structure[1])
	}
	l.addToProgram(reg)

	// String literal
	str := strings.Join(structure[2:], " ")
	str = str[1:]
	str = str[:len(str)-1]

	strLen := len(str)
	if strLen > 32768 {
		return fmt.Errorf("String too long on line %d", l.line)
	}

	// Push string lenth
	l.addToProgram(byte(strLen >> 8))
	l.addToProgram(byte(strLen))

	// Add string literal
	l.addSliceToProgram([]byte(str))
	return nil
}

func (l *Lexer) parseParamsIntReg(structure []string) error {
	if len(structure) != 3 {
		return fmt.Errorf("Expected register and int on line %d", l.line)
	}

	// Value
	code, err := strconv.ParseInt(structure[1], 10, 64)
	if err != nil {
		return err
	}
	l.addSliceToProgram(intToBytes(code))

	if structure[2][0] != '$' {
		return fmt.Errorf("Expected register on line %d", l.line)
	}

	// Register
	reg, ok := getRegister(structure[2][1:])
	if !ok {
		return fmt.Errorf("%s is not a register", structure[1])
	}
	l.addToProgram(reg)

	return nil
}

func (l *Lexer) subLabels() error {
	for _, sub := range l.labelSubs {
		loc, ok := l.labels[sub.label]
		if !ok {
			return fmt.Errorf("Label %s not defined", sub.label)
		}

		locBytes := intToBytes(loc)
		l.program[sub.pos] = locBytes[0]
		l.program[sub.pos+1] = locBytes[1]
		l.program[sub.pos+2] = locBytes[2]
		l.program[sub.pos+3] = locBytes[3]
		l.program[sub.pos+4] = locBytes[4]
		l.program[sub.pos+5] = locBytes[5]
		l.program[sub.pos+6] = locBytes[6]
		l.program[sub.pos+7] = locBytes[7]
	}
	return nil
}

func intToBytes(i int64) []byte {
	out := make([]byte, 8)
	binary.PutVarint(out, i)
	return out
}

func getRegister(name string) (byte, bool) {
	v, ok := registers[strings.ToUpper(name)]
	return v, ok
}
