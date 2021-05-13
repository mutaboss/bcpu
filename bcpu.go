package main

const DefaultMemorySize = 4096
const ProgramStart = 256
const InstructionWidth = 2

const (
    OpHalt = iota
    OpNoop
)

type Bcpu struct {
    pc int
    memory [DefaultMemorySize]int16
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

func (cpu *Bcpu) SetMemory(location int, value int16) {
    cpu.memory[location] = value
}

func (cpu *Bcpu) Run() {
    cpu.pc = ProgramStart
    for {
        inst, _ := cpu.memory[cpu.pc], cpu.memory[cpu.pc+1]
        cpu.pc = cpu.pc + 2
        if inst == OpHalt {
            break
        } else if inst == OpNoop {
            /* do nothing */
        }
    }
}

