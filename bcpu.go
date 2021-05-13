package main

import (
    "fmt"
)

const DefaultMemorySize = 4096
const ProgramStart = 256
const InstructionWidth = 2
const RegisterCount = 16

const (
    OpHalt = iota
    OpNoop
    OpSetReg
)

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
        inst, param := cpu.memory[cpu.pc], cpu.memory[cpu.pc+1]
        opcode := (inst & 0xFF00) >> 8
        target := inst&0x00FF
        cpu.pc = cpu.pc + InstructionWidth
        if opcode == OpHalt {
            break
        } else if opcode == OpNoop {
            /* do nothing */
        } else if opcode == OpSetReg {
            if target < 0 || target > RegisterCount - 1 {
                return fmt.Errorf("Invalid register %d.", target)
            }
            cpu.register[target] = param
        } else {
            return fmt.Errorf("Invalid opcode: %d.", opcode)
        }
    }
    return nil
}
