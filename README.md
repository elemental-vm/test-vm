# TestVM

**This project is a test bed and thus has no name. It shall be referred to as TestVM.**

TestVM is an exploratory project to learn about bytecode virtual machines and how they're implemented.
TestVM is a hybrid stack/register based virtual machine. A stack is used for function calls while registers
are used for function parameters and temporary values.

Programs are written in plain text. Binary files are not supported yet.

## Instructions

The instructions set is quite limited, but provides enough for simple programs.

In the following table, `#` refers to any 64bit signed integer. `$reg` refers to a register named A-J, RT, PC, SP, or FP.
`%label` refers to a label in code, explanation below. `TOS` refers to the top of stack.

| Code | Name    | Syntax              | Desc.                                                                          |
|------|---------|---------------------|--------------------------------------------------------------------------------|
| 0x00 | HALT    | HALT # EXIT#        | Stop program execution returning exit code #. Exit code must be between 0-255. |
| 0x01 | PUSHI   | PUSHI 42            | Push an integer onto the stack.                                                |
| 0x02 | PUSHSTR | PUSHSTR "Hello"     | Push a string onto the stack.                                                  |
| 0x03 | PUSHREG | PUSHREG $reg        | Push the value in $reg onto the stack.                                         |
| 0x04 | POP     | POP                 | Pop TOS. Discards value.                                                       |
| 0x05 | POPREG  | POPREG $reg         | Pop TOS value into $reg.                                                       |
| 0x06 | STORE   | STORE $reg          | Store TOS value to $reg without popping the stack.                             |
| 0x07 | SWAP    | SWAP                | Swap the two TOS values. E.g: [1, 3, 4, 7] -> [3, 1, 4, 7].                    |
| 0x08 | DUP     | DUP                 | Push a copy of TOS onto the stack. E.g: [3, 4, 7] -> [3, 3, 4, 7].             |
| 0x09 | ADD     | ADD                 | Add the two TOS values, pushes result onto stack.                              |
| 0x0A | SUB     | SUB                 | Subtract the two TOS values, pushes result onto stack.                         |
| 0x0B | MUL     | MUL                 | Multiply the two TOS values, pushes result onto stack.                         |
| 0x0C | DIV     | DIV                 | Divide the two TOS values, pushes result onto stack.                           |
| 0x0D | SETI    | SETI $reg #/%label  | Set $reg to # or the memory location of %label.                                |
| 0x0E | SETSTR  | SETSTR $reg "Hello" | Set $reg to string.                                                            |
| 0x0F | JUMP    | JMP #/%label        | Unconditionally jump to location.                                              |
| 0x10 | JUMPGTZ | JMPGTZ #/%label     | Jump to location if TOS is greater than 0.                                     |
| 0x11 | JUMPLTZ | JMPLTZ #/%label     | Jump to location if TOS is less than 0.                                        |
| 0x12 | JUMPEQ  | JMPEQ #/%label      | Jump to location if TOS is equal to 0.                                         |
| 0x13 | JUMPNEQ | JMPNEQ #/%label     | Jump to location if TOS is not equal to 0.                                     |
| 0x14 | PRINT   | PRINT               | Print TOS value to stdout.                                                     |
| 0x15 | PRINTR  | PRINTR $reg         | Print value of $reg.                                                           |
| 0x16 | DUMP    | DUMP                | Print the full stack to stdout.                                                |
| 0x17 | DUMPR   | DUMPR               | Print all registers to stdout.                                                 |
| 0x18 | RETURN  | RETURN              | Return to address in $RT                                                       |
| 0x19 | CALL    | CALL #/%label       | Call location as a function, stores return address in $RT                      |
| 0x1A | CONCAT  | CONCAT              | Concatenate the top two stack values. Places result on TOS.                    |
| 0x1B | PARAM   | PARAM $reg #        | Move parameter # to $reg.                                                      |
| 0x1C | JUMPREG | JMPREG $reg         | Jump to location store in $reg.                                                |
| 0x1D | COMPARE | CMP $reg $reg       | Compare the values of two registers. Sets the zero flag.                       |
| 0x1E | JUMPZGTZ| JMPZGTZ #/%label    | Jump to location if the zero flag is greater than 0.                           |
| 0x1F | JUMPZLTZ| JMPZLTZ #/%label    | Jump to location if the zero flag is less than 0.                              |
| 0x20 | JUMPZEQ | JMPZEQ #/%label     | Jump to location if the zero flag is equal to 0.                               |
| 0x21 | JUMPZNEQ| JMPZNEQ #/%label    | Jump to location if the zero flag is not equal to 0.                           |

## Labels

Labels can be used in source code in place for instruction numbers for locations. It's highly recommended to use labels
as they will automatically adjust with code changes while hard-coded addresses won't. A line can have a label applied to
it by prefixing the line with characters followed by a colon. Labels can be referenced by prefixing a label with a percent
sign. Labels may be used before they're defined. The locations are inserted after the full source has been parsed.

```asm
setup:  SETI $A 0
        SETI $B 1
loop:   PUSHREG $A
        PUSHREG $B
        ADD
        JMPLZ %exit
        PRINT
        DUP
        POPREG $A
        POPREG $B
        JMP %loop
exit:   HALT 0
```

## Comments

Lines beginning with a semicolon are considered comments. A comment may begin anywhere in a line and will continue to the
end of the line. There are no block comments.

```asm
;; This is a comment for my awesome program

        PUSHI 0     ; Push a seed value
loop:   PUSHI 1
        ADD         ; Add one
        PRINT       ; Print value
        JMP %loop   ; Loop indefinitely
```

## Registers

TestVM has 10 general purpose registers and four special purpose registers. Registers 'A' through 'J' may be used however the
programmer likes. The programmer is responsible for preserving them between function calls.

In source code, registers are denoted with a dollar sign: `$A`.

Registers PC, SP, FP, and RT are special purpose. The registers are for the Program Counter, Stack Pointer, Frame Pointer, and
Return address respectively. Currently offsets are not possible and the FP and RT registers aren't used. They will eventually
be used once proper function calls are implemented.
