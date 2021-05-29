package main

import (
	"fmt"
    "unicode"
)

const DefaultMemorySize uint16 = 4096
const ProgramStart uint16 = 256
const RegisterCount uint16 = 16

// ************************************************************************************************
// Opcodes

type Opcode uint16

const (
	OpHalt   Opcode = 0    // HLT
	OpNoop          = 1    // NOP
	OpJmp           = 2    // JMP ADDR
	OpJeq           = 3    // JEQ ADDR
	OpJgt           = 4    // JGT ADDR
	OpJlt           = 5    // JLT ADDR
	OpSetReg        = 8    // STR REG, LITERAL
	OpLoad          = 9    // LDR REG, ADDR
	OpStore         = 10   // STR REG, ADDR
	OpAddReg        = 11   // ADD SRC DST
	OpSubReg        = 12   // SUB SRC DST
	OpMulReg        = 13   // MUL SRC DST
	OpDivReg        = 14   // DIV SRC DST
	OpCmp           = 15   // CMP SRC DST
    OpAnd           = 16   // AND SRC DST
    OpOr            = 17   // OR  SRC DST
    OpXor           = 18   // XOR SRC DST
    OpShl           = 19   // SHL AMT DST
    OpShr           = 20   // SHR AMT DST
    OpNot           = 21   // NOT DST
)

// ************************************************************************************************
// Instructions

type Instruction struct {
	opcode Opcode
	regsrc uint16
	regtgt uint16
	memloc uint16
}

func msk(num uint16, mask uint16, shift uint16) uint16 {
	return (num & mask) << shift
}

func umsk(num uint16, mask uint16, shift int) uint16 {
	return (num & mask) >> shift
}

/*
 * Opcode with Embedded memory reference:
 * 0aaammmmmmmmmmmm
 *
 * Opcode with registers:
 * 1aaaaaaxsssstttt
 *
 * a = opcode
 * m = memory address (12 bit = 4096)
 * s = source register
 * t = target register
 */

func NewInstruction(op Opcode, regsrc uint16, regtgt uint16, memloc uint16) *Instruction {
	inst := new(Instruction)
	inst.opcode = op
	if op < 8 {
		inst.memloc = memloc & 0x0fff
	} else {
		inst.regsrc = regsrc & 0x000f
		inst.regtgt = regtgt & 0x000f
	}
	return inst
}

func DecodeInstruction(instruction uint16) *Instruction {
	inst := new(Instruction)
	if instruction&0x8000 == 0 { // iiiimmmmmmmmmmmm
		inst.opcode = Opcode(umsk(instruction, 0x7000, 12))
		inst.memloc = umsk(instruction, 0x0fff, 0)
	} else { // 1iiiiii0sssstttt
		inst.opcode = Opcode(umsk(instruction, 0x7e00, 9)) + 7
		inst.regsrc = umsk(instruction, 0x00f0, 4)
		inst.regtgt = umsk(instruction, 0x000f, 0)
	}
	return inst
}

func (instruction *Instruction) Encode() uint16 {
	if instruction.opcode < 8 {
		return msk(uint16(instruction.opcode), 0x0007, 12) +
			msk(instruction.memloc, 0x0fff, 0)
	} else {
		return 0x8000 +
			msk(uint16(instruction.opcode-7), 0x003f, 9) +
			msk(instruction.regsrc, 0x000f, 4) +
			msk(instruction.regtgt, 0x000f, 0)
	}
}

// ************************************************************************************************
// Bcpu: The Machine.

type Bcpu struct {
	pc       uint16
	memory   [DefaultMemorySize]uint16
	register [RegisterCount]uint16
	flags    uint8
}

func NewBcpu() *Bcpu {
	cpu := new(Bcpu)
	cpu.pc = ProgramStart
	return cpu
}

// TODO: Alter jump addresses based on memstart.
func (cpu *Bcpu) Load(memstart uint16, instructions []uint16) {
    for i := 0; i < len(instructions); i++ {
        cpu.SetMemory(memstart + uint16(i), instructions[i])
    }
}

func (cpu *Bcpu) setOverflow() {
	cpu.flags |= 0b00000001
}

func (cpu *Bcpu) unsetOverflow() {
	cpu.flags &= 0b11111110
}

func (cpu *Bcpu) GetOverflow() bool {
	return cpu.flags&0b00000001 == 1
}

func (cpu *Bcpu) setEqual() {
	cpu.flags &= 0b11111001
}

func (cpu *Bcpu) setGreater() {
	cpu.setEqual()
	cpu.flags |= 0b11111101
}

func (cpu *Bcpu) setLesser() {
	cpu.setEqual()
	cpu.flags |= 0b11111011
}

func (cpu *Bcpu) GetEqual() bool {
	return cpu.flags&0b00000110 == 0
}

func (cpu *Bcpu) GetGreater() bool {
	return cpu.flags&0b00000100 > 0
}

func (cpu *Bcpu) GetLesser() bool {
	return cpu.flags&0b00000010 > 0
}

func (cpu *Bcpu) MemorySize() uint16 {
	return uint16(len(cpu.memory))
}

func (cpu *Bcpu) ProgramCounter() uint16 {
	return cpu.pc
}

func (cpu *Bcpu) SetMemory(location uint16, value uint16) error {
	if location >= DefaultMemorySize {
		return fmt.Errorf("Invalid memory location: %d.", location)
	}
	cpu.memory[location] = value
	return nil
}

func (cpu *Bcpu) GetMemory(location uint16) (uint16, error) {
	if location >= DefaultMemorySize {
		return 0, fmt.Errorf("Invalid memory location: %d.", location)
	} else {
		return cpu.memory[location], nil
	}
}

func (cpu *Bcpu) SetRegister(reg uint16, val uint16) error {
	if reg > RegisterCount-1 {
		return fmt.Errorf("Bad register designation: %d.", reg)
	}
	cpu.register[reg] = val
	return nil
}

func (cpu *Bcpu) GetRegister(reg uint16) (uint16, error) {
	if reg > RegisterCount-1 {
		return 0, fmt.Errorf("Bad register designation: %d.", reg)
	}
	return cpu.register[reg], nil
}

func (cpu *Bcpu) Run() error {
	cpu.pc = ProgramStart
	for keep_going := true; keep_going; {
		inst := DecodeInstruction(cpu.memory[cpu.pc])
		cpu.pc++
		switch inst.opcode {
		case OpHalt:
			keep_going = false
		case OpNoop:
			/* do nothing */
		case OpSetReg:
			param := cpu.memory[cpu.pc]
			cpu.pc++
			cpu.register[inst.regtgt] = param
		case OpLoad:
			location := cpu.memory[cpu.pc]
			cpu.pc++
			cpu.register[inst.regtgt] = cpu.memory[location]
		case OpStore:
			location := cpu.memory[cpu.pc]
			cpu.pc++
			cpu.memory[location] = cpu.register[inst.regsrc]
		case OpAddReg:
			newval := int32(cpu.register[inst.regsrc]) + int32(cpu.register[inst.regtgt])
			cpu.register[inst.regtgt] = uint16(newval)
			if newval < 0 || newval > 65535 {
				cpu.setOverflow()
			} else {
				cpu.unsetOverflow()
			}
		case OpSubReg:
			newval := int32(cpu.register[inst.regsrc]) - int32(cpu.register[inst.regtgt])
			cpu.register[inst.regtgt] = uint16(newval)
			if newval < 0 || newval > 65535 {
				cpu.setOverflow()
			} else {
				cpu.unsetOverflow()
			}
		case OpMulReg:
			newval := int32(cpu.register[inst.regsrc]) * int32(cpu.register[inst.regtgt])
			cpu.register[inst.regtgt] = uint16(newval)
			if newval < 0 || newval > 65535 {
				cpu.setOverflow()
			} else {
				cpu.unsetOverflow()
			}
		case OpDivReg:
			newval := int32(cpu.register[inst.regsrc]) / int32(cpu.register[inst.regtgt])
			cpu.register[inst.regtgt] = uint16(newval)
			if newval < 0 || newval > 65535 {
				cpu.setOverflow()
			} else {
				cpu.unsetOverflow()
			}
		case OpJmp:
			cpu.pc = inst.memloc
		case OpCmp:
			a, b := cpu.register[inst.regsrc], cpu.register[inst.regtgt]
			if a == b {
				cpu.setEqual()
			} else if a > b {
				cpu.setGreater()
			} else {
				cpu.setLesser()
			}
		case OpJeq:
			if cpu.GetEqual() {
				cpu.pc = inst.memloc
			}
		case OpJgt:
			if cpu.GetGreater() {
				cpu.pc = inst.memloc
			}
		case OpJlt:
			if cpu.GetLesser() {
				cpu.pc = inst.memloc
			}
        case OpShl:
            cpu.register[inst.regtgt] = cpu.register[inst.regtgt] << inst.regsrc
        case OpShr:
            cpu.register[inst.regtgt] = cpu.register[inst.regtgt] >> inst.regsrc
        case OpAnd:
            cpu.register[inst.regtgt] = cpu.register[inst.regsrc] & cpu.register[inst.regtgt]
        case OpOr:
            cpu.register[inst.regtgt] = cpu.register[inst.regsrc] | cpu.register[inst.regtgt]
        case OpXor:
            cpu.register[inst.regtgt] = cpu.register[inst.regsrc] ^ cpu.register[inst.regtgt]
        case OpNot:
            cpu.register[inst.regtgt] = ^ cpu.register[inst.regtgt]
		default:
			return fmt.Errorf("Invalid opcode: %d (%b).", inst.opcode, inst.Encode())
		}
	}
	return nil
}

// ************************************************************************************************
// Lexer

type TokenType uint16

const (
	TokHalt   TokenType = iota
	TokNoop
	TokJmp
	TokJeq
	TokJgt
	TokJlt
	TokSetReg
	TokLoad
	TokStore
	TokAddReg
	TokSubReg
	TokMulReg
	TokDivReg
	TokCmp
    TokAnd
    TokOr
    TokXor
    TokShl
    TokShr
    TokNot
    TokInteger
    TokEOF
)

type Token struct {
    literal string
    token TokenType
}

type Lexer struct {
    input []rune
    position int
    readPosition int
    ch    rune
    line int
}

func (l *Lexer) NewLexer(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) readChar() {
    if l.ch == '\n' {
        l.line += 1
    }
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
        l.position = l.readPosition
        l.readPosition += 1
    }
}

func (l *Lexer) NextToken() *Token {
    var tok Token
    if l.ch == 0 {
        return &Token{"EOF", TokEOF}
    } else if unicode.IsDigit(l.ch) {
        return l.readNumber()
    } else if unicode.IsLetter(l.ch) {
        return l.readIdentifier()
    } else {
        fmt.Printf("Invalid character line %d: %s.\n", l.line, l.ch)
        l.readChar()
    }
    return nil // We should not be here!
}

func (l *Lexer) readNumber() *Token {
    var result []rune
    for unicode.IsDigit(l.ch) {
        result := append(result, l.ch)
        l.readChar()
    }
    return &Token{string(result), TokNumber}
}

func (l *Lexer) readIdentifier() *Token {
    var result []rune
    for unicode.IsLetter(l.ch) {
        result := append(result, l.ch)
        l.readChar()
    }
    if ident, err := l.lookupIdentifier(string(result)); err != nil {
        fmt.
    return &Token{string(result), l.lookupIdentifier(string(result))}
}

func (l *Lexer) lookupIdentifier(Identifier) TokenType {
}

// ************************************************************************************************
// Assembler

type Assembler struct {
    input string
    runes []rune
    index int
}

// TODO: Flesh out the lexer. Should it be a separate object?

func (asm *Assembler) getNextChar() rune

func (asm *Assembler) getNextToken() (*Token) {
    if asm.index >= len(asm.runes) {
        return Token{"EOF", TokEOF}
    }
    ch := asm.runes[asm.index++]
    if unicode.IsDigit(ch) {
        // parse a number
        tokstr := []runes{ch}
        
    i := asm.index
    for _, v := range asm.input[i:] {
        if !unicode.IsSpace(v) {
            break
        }
        i++
    }
    if i >= len(asm.input) {
        return "", true
    }
    j := i
    for _, v := range(asm.input[j:]) {
        if unicode.IsSpace(v) {
            break
        }
        j++
    }
    asm.index = j
    return asm.input[i:j], false
}

// TODO: Do this.
func Assemble(input string) ([]uint16, error) {
    asm := Assembler{input, []rune(input), 0}
    var result []uint16
    for tok, eof := asm.getNextToken(); !eof; tok, eof = asm.getNextToken() {
        // switch tok {
        // case "SETR":
        //     // TODO: Fill this in.
    }
    return result, nil
}
