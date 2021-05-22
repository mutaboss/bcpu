package main

import (
    "fmt"
)

const DefaultMemorySize = 4096
const ProgramStart = 256
const RegisterCount = 16

// ************************************************************************************************
// Opcodes

type Opcode uint16

const (
    OpHalt   Opcode = 0
    OpNoop          = 1
    OpSetReg        = 8
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
 * 1aaaaxxxsssstttt
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
    if instruction & 0x8000 == 0 { // iiiimmmmmmmmmmmm
        inst.opcode = Opcode(umsk(instruction,0x7000,12))
        inst.memloc = umsk(instruction,0x0fff,0)
    } else { // 1iiii000sssstttt
        inst.opcode = Opcode(umsk(instruction,0x7800,11)) + 7
        inst.regsrc = umsk(instruction,0x00f0,4)
        inst.regtgt = umsk(instruction,0x000f,0)
    }
    return inst
}

func (instruction *Instruction) Encode() uint16 {
    if instruction.opcode < 8 {
        return msk(uint16(instruction.opcode),0x0007,12) +
            msk(instruction.memloc,0x0fff,0)
    } else {
        return 0x8000 +
            msk(uint16(instruction.opcode-7),0x000f,11) +
            msk(instruction.regsrc,0x000f,4) +
            msk(instruction.regtgt,0x000f,0)
    }
}

func EncodeOpcode(op Opcode, regsrc uint16, regtgt uint16, memloc uint16) (uint16, error) {
    if op < 8 {
        return msk(uint16(op),0x0007,8) + msk(memloc,0x0fff,0), nil
    } else {
        return 0x8000 + msk(uint16(op),0x000f,7) + msk(regsrc,0x000f,4) + msk(regtgt,0x000f,0), nil
    }
}

func DecodeOpcode(instruction uint16) (Opcode, uint16, uint16, uint16, error) {
    if instruction & 0x8000 == 0x8000 {
        return Opcode(umsk(instruction,0x7000,8)), 0, 0, umsk(instruction,0x0fff,0), nil
    } else {
        return Opcode(umsk(instruction,0x7800,7)), umsk(instruction,0x00f0,4), umsk(instruction,0x000f,0), 0, nil
    }
}

// ************************************************************************************************
// Bcpu: The Machine.

type Bcpu struct {
    pc int
    memory [DefaultMemorySize]uint16
    register [RegisterCount]uint16
}

func NewBcpu() *Bcpu {
    cpu := new(Bcpu)
    cpu.pc = ProgramStart
    return cpu
}

func (cpu *Bcpu) MemorySize() int {
    return len(cpu.memory)
}

func (cpu *Bcpu) ProgramCounter() int {
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

func (cpu *Bcpu) GetRegister(reg int) (uint16, error) {
    if reg < 0 || reg > RegisterCount - 1 {
        return 0, fmt.Errorf("Bad register designation: %d.", reg)
    }
    return cpu.register[reg], nil
}

func (cpu *Bcpu) Run() error {
    cpu.pc = ProgramStart
    for {
        inst := DecodeInstruction(cpu.memory[cpu.pc])
        cpu.pc ++
        if inst.opcode == OpHalt {
            break
        } else if inst.opcode == OpNoop {
            /* do nothing */
        } else if inst.opcode == OpSetReg {
            param := cpu.memory[cpu.pc]
            cpu.pc ++
            cpu.register[inst.regtgt] = param
        } else {
            return fmt.Errorf("Invalid opcode: %d (%b).", inst.opcode, inst.Encode())
        }
    }
    return nil
}
