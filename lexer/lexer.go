package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/elemental-vm/test-vm/vm"
)

type Lexer struct {
	r *bufio.Reader

	program []int64
}

func New(file string) (*Lexer, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		r: bufio.NewReader(f),
	}, nil
}

func (l *Lexer) addToProgram(word int64) {
	l.program = append(l.program, word)
}

func (l *Lexer) Parse() ([]int64, error) {
	l.program = make([]int64, 0, 500)
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
		line = strings.TrimSpace(line)
		if line == "" || line[0] == ';' {
			continue
		}

		text := strings.Split(line, ";")
		line = strings.TrimSpace(text[0])

		structure := strings.Split(line, " ")
		parts := len(structure)
		if parts == 0 {
			continue
		}

		bytecode, ok := bytecodes[structure[0]]
		if !ok {
			return nil, fmt.Errorf("Unknown code %s", structure[0])
		}
		l.addToProgram(bytecode)

		switch bytecode {
		case vm.Halt:
			if parts != 2 {
				return nil, errors.New("Expected int after HALT")
			}
			code, err := strconv.ParseInt(structure[1], 10, 64)
			if err != nil {
				return nil, err
			}
			l.addToProgram(code)
		case vm.Push:
			if parts != 2 {
				return nil, errors.New("Expected int after PUSH")
			}
			code, err := strconv.ParseInt(structure[1], 10, 64)
			if err != nil {
				return nil, err
			}
			l.addToProgram(code)
		case vm.PushReg:
			if parts != 2 {
				return nil, errors.New("Expected register after PUSHREG")
			}
			reg, ok := registers[structure[1]]
			if !ok {
				return nil, fmt.Errorf("%s is not a register", structure[1])
			}
			l.addToProgram(reg)
		case vm.PopReg:
			if parts != 2 {
				return nil, errors.New("Expected register after POPREG")
			}
			reg, ok := registers[structure[1]]
			if !ok {
				return nil, fmt.Errorf("%s is not a register", structure[1])
			}
			l.addToProgram(reg)
		case vm.Set:
			if parts != 3 {
				return nil, errors.New("Expected register and int after SET")
			}
			// Register
			reg, ok := registers[structure[1]]
			if !ok {
				return nil, fmt.Errorf("%s is not a register", structure[1])
			}
			l.addToProgram(reg)

			// Value
			code, err := strconv.ParseInt(structure[2], 10, 64)
			if err != nil {
				return nil, err
			}
			l.addToProgram(code)
		case vm.Jump:
			if parts != 2 {
				return nil, errors.New("Expected int after JMP")
			}
			code, err := strconv.ParseInt(structure[1], 10, 64)
			if err != nil {
				return nil, err
			}
			l.addToProgram(code)
		case vm.JumpGtz:
			if parts != 2 {
				return nil, errors.New("Expected int after JMPGZ")
			}
			code, err := strconv.ParseInt(structure[1], 10, 64)
			if err != nil {
				return nil, err
			}
			l.addToProgram(code)
		}
	}

	return l.program, nil
}
