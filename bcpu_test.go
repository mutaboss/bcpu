package main

import (
    "testing"
    "fmt"
)

func TestBcpuNew(t *testing.T) {
    var cpu *Bcpu = NewBcpu()
    if cpu.MemorySize() != DefaultMemorySize {
        t.Error("Memory is not correctly allocated.")
    }
    if cpu.ProgramCounter() != ProgramStart {
        t.Error("Program Counter is not {}", ProgramStart)
    }
}

func TestBcpuRun(t *testing.T) {
    var cpu *Bcpu = NewBcpu()
    cpu.Run()
    if cpu.ProgramCounter() != ProgramStart + InstructionWidth {
        t.Error(fmt.Sprintf("Program Counter did not advance (is %d).", cpu.ProgramCounter()))
    }
}

func TestBcpuNoop(t *testing.T) {
    var cpu *Bcpu = NewBcpu()
    cpu.SetMemory(ProgramStart, OpNoop)
    cpu.SetMemory(ProgramStart+2, OpNoop)
    cpu.Run()
    if cpu.ProgramCounter() != ProgramStart + 3*InstructionWidth {
        t.Error(fmt.Sprintf("Program Counter should be %d, is %d.", ProgramStart + 2, cpu.ProgramCounter()))
    }
}
