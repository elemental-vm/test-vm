# TestVM

**This project is a test bed and thus has no name. It shall be referred to as TestVM.**

TestVM is an exploratory project to learn about bytecode virtual machines and how they're implemented.
TestVM is a hybrid stack/register based virtual machine. A stack is used for function calls while registers
are used for function parameters and temporary values.

Programs are written in plain text. Binary files are not supported yet.

## Instructions

The instructions set is quite limited, but provides enough for simple programs.

In the following list, `#` refers to any 64bit signed integer. `reg` refers to a register name A-J, RT.
`%label` refers to a label in code, explanation below. `TOS` refers to the top of stack.

0. `EXIT/HALT #`: Stop execution and exit with code #.
1. `PUSH #`: Push a value onto the stack.
2. `PUSHREG reg`: Push a value from a register onto the stack.
3. `POP`: Pop the TOS value, discard value.
4. `POPREG reg`: Pop the TOS value into a register.
5. `SWAP`: Swap the two TOS values. E.g: [1, 3, 4, 7] -> [3, 1, 4, 7].
6. `DUP`: Push a copy of TOS onto the stack. E.g: [3, 4, 7] -> [3, 3, 4, 7].
7. `SET reg #`: Set the value of register.
8. `ADD`: Add the two TOS values, pushes result onto stack.
9. `SUB`: Subtract the two TOS values, pushes result onto stack.
10. `MUL`: Multiply the two TOS values, pushes result onto stack.
11. `DIV`: Divide the two TOS values, pushes result onto stack.
12. `JMP #/%label`: Unconditionally jump to location.
13. `JMPGZ #/%label`: Jump to location if TOS is greater than 0.
14. `JMPLZ #/%label`: Jump to location if TOS is less than 0.
15. `JMPEQ #/%label`: Jump to location if TOS is equal to 0.
16. `JMPNEQ #/%label`: Jump to location if TOS is not equal to 0.
17. `PRINT`: Print the TOS value to stdout.
18. `PRINTS`: Print the stack to stdout.
19. `CALL #/%label`: Call location as a function, pushes the return address to the stack. *
20. `RET`: Return from a function call to the callee. *

\* Function calling is not yet finalized. It needs work.

## Labels

Labels can be used in source code in place for instruction numbers for locations. It's highly recommended to use labels
as they will automatically adjust with code changes while hard-coded addresses won't. A line can have a label applied to
it by prefixing the line with characters followed by a colon. Labels can be referenced by prefixing a label with a percent
sign. Labels may be used before they're defined. The locations are inserted after the full source has been parsed.

```asm
setup:  SET A 0
        SET B 0
loop:   PUSHREG A
        PUSHREG B
        ADD
        PRINT
        DUP
        POPREG A
        POPREG B
        JMP %loop
```

## Comments

Lines beginning with a semicolon are considered comments. A comment may begin anywhere in a line and will continue to the
end of the line. There are no block comments.

```asm
;; This is a comment for my awesome program

        PUSH 0      ; Push a seed value
loop:   PUSH 1
        ADD         ; Add one
        PRINT       ; Print value
        JMP %loop   ; Loop indefinitely
```

## Registers

TestVM has 10 general purpose registers and one reserved register. Registers 'A' through 'J' may be used however the programmer
likes. The programmer is responsible for preserving them between function calls.

Register 'RT' is reserved. I intend to use it for function calls eventually. Use it at your own risk.
