package main

import (
    "fmt"
)

const DefaultMemorySize = 4096
const ProgramStart = 256
const RegisterCount = 16

type Opcode uint16

const (
    OpHalt   Opcode = 0
    OpNoop          = 1
    OpSetReg        = 2
)

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

func NewInstruction(instruction uint16) *Instruction {
    inst := new(Instruction)
    if instruction & 0x8000 == 0x8000 {
        inst.opcode = Opcode(umsk(instruction,0x7000,8))
        inst.memloc = umsk(instruction,0x0fff,0)
    } else {
        inst.opcode = Opcode(umsk(instruction,0x7800,7))
        inst.regsrc = umsk(instruction,0x00f0,4)
        inst.regtgt = umsk(instruction,0x000f,0)
    }
    return inst
}

func (instruction *Instruction) Encode() uint16 {
    if instruction.opcode < 8 {
        return msk(uint16(instruction.opcode),0x0007,8) +
            msk(instruction.memloc,0x0fff,0)
    } else {
        return 0x8000 +
            msk(uint16(instruction.opcode),0x000f,7) +
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

func (cpu *Bcpu) SetMemory(location int, value uint16) {
    cpu.memory[location] = value
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
        inst := NewInstruction(cpu.memory[cpu.pc])
        cpu.pc ++
        var param uint16 = 0
        if inst.opcode > 8 {
            param = cpu.memory[cpu.pc]
            cpu.pc ++
        }
        if inst.opcode == OpHalt {
            break
        } else if inst.opcode == OpNoop {
            /* do nothing */
        } else if inst.opcode == OpSetReg {
            cpu.register[inst.regtgt] = param
        } else {
            return fmt.Errorf("Invalid opcode: %d.", inst.opcode)
        }
    }
    return nil
}
