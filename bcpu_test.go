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
    if err := cpu.Run(); err != nil {
        t.Error(fmt.Sprintf("Execution error: %s.", err))
    }
    if cpu.ProgramCounter() != ProgramStart + 1 {
        t.Error(fmt.Sprintf("Program Counter did not advance (is %d).", cpu.ProgramCounter()))
    }
}

func TestBcpuNoop(t *testing.T) {
    var cpu *Bcpu = NewBcpu()
    cpu.SetMemory(ProgramStart, OpNoop<<8)
    cpu.SetMemory(ProgramStart+1, OpNoop<<8)
    if err := cpu.Run(); err != nil {
        t.Error(fmt.Sprintf("Execution error: %s.", err))
    }
    exppc := ProgramStart + 3
    if cpu.ProgramCounter() != exppc {
        t.Error(fmt.Sprintf("Program Counter should be %d, is %d.", exppc, cpu.ProgramCounter()))
    }
}

func TestBadOpcode(t *testing.T) {
    var cpu *Bcpu = NewBcpu()
    cpu.SetMemory(ProgramStart, 0xffff)
    if err := cpu.Run(); err == nil {
        t.Error("Expected 0xFFFF to be a bad opcode.")
    }
}

func testRegisterGet(cpu *Bcpu, reg int, expval uint16) bool {
    val, err := cpu.GetRegister(reg)
    if err != nil {
        return false
    }
    if val != expval {
        return false
    }
    return true
}

func testRegisterSet(t *testing.T, cpu *Bcpu, reg uint16, expval uint16) {
    cpu.SetMemory(ProgramStart, OpSetReg<<8 + reg)
    cpu.SetMemory(ProgramStart+1, expval)
    if err := cpu.Run(); err != nil {
        t.Error(fmt.Sprintf("Execution error: %s", err))
    }
    val, err := cpu.GetRegister(int(reg))
    if err != nil {
        t.Error(err)
    }
    if val != expval {
        t.Error(fmt.Sprintf("Expected register %d to be %d, was %d.", reg, expval, val))
    }
}

func TestBcpuRegisters(t *testing.T) {
    var cpu *Bcpu = NewBcpu()
    if ! testRegisterGet(cpu, 0, 0) {
        t.Error(fmt.Sprintf("Expected register %d to be good, and have a value of %d.", 0, 0))
    }
    if ! testRegisterGet(cpu, RegisterCount-1, 0) {
        t.Error(fmt.Sprintf("Expected register %d to be good, and have a value of %d.", RegisterCount-1, 0))
    }
}

func TestBcpuOpSetreg(t *testing.T) {
    var cpu *Bcpu = NewBcpu()
    testRegisterSet(t, cpu, 0, 16)
    testRegisterSet(t, cpu, 1, 256)
    testRegisterSet(t, cpu, RegisterCount-1, 256)
}

// func TestBcpuEncodeOpcode(t *testing.T) {
//     if opcode, err := EncodeOpcode(OpHalt, 0, 0, 0); err != nil {
//         t.Error("Invalid opcode conversion.")
//     } else if opcode != 0 {
//         t.Error("Halt instruction should give a 0 opcode.")
//     }
// }
