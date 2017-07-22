package lexer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/elemental-vm/test-vm/vm"
)

type Lexer struct {
	r    *bufio.Reader
	line int

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
		l.line++

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
			return nil, fmt.Errorf("Unknown instruction %s on line %d", structure[0], l.line)
		}
		l.addToProgram(bytecode)

		switch bytecode {
		case vm.Halt:
			err = l.parseParamOneInt(structure)
		case vm.Push:
			err = l.parseParamOneInt(structure)
		case vm.PushReg:
			err = l.parseParamOneRegister(structure)
		case vm.PopReg:
			err = l.parseParamOneRegister(structure)
		case vm.Set:
			err = l.parseParamsRegInt(structure)
		case vm.Jump:
			err = l.parseParamOneInt(structure)
		case vm.JumpGtz:
			err = l.parseParamOneInt(structure)
		default:
			err = nil
		}

		if err != nil {
			return nil, err
		}
	}

	return l.program, nil
}

func (l *Lexer) parseParamOneInt(structure []string) error {
	if len(structure) != 2 {
		return fmt.Errorf("Expected int on line %d", l.line)
	}
	code, err := strconv.ParseInt(structure[1], 10, 64)
	if err != nil {
		return err
	}
	l.addToProgram(code)
	return nil
}

func (l *Lexer) parseParamOneRegister(structure []string) error {
	if len(structure) != 2 {
		return fmt.Errorf("Expected register on line %d", l.line)
	}
	reg, ok := registers[structure[1]]
	if !ok {
		return fmt.Errorf("%s is not a register", structure[1])
	}
	l.addToProgram(reg)
	return nil
}

func (l *Lexer) parseParamsRegInt(structure []string) error {
	if len(structure) != 3 {
		return fmt.Errorf("Expected register and int on line %d", l.line)
	}

	// Register
	reg, ok := registers[structure[1]]
	if !ok {
		return fmt.Errorf("%s is not a register", structure[1])
	}
	l.addToProgram(reg)

	// Value
	code, err := strconv.ParseInt(structure[2], 10, 64)
	if err != nil {
		return err
	}
	l.addToProgram(code)
	return nil
}
