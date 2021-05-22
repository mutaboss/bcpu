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

func TestBcpuMemory(t *testing.T) {
    var cpu *Bcpu = NewBcpu()
    if err := cpu.SetMemory(ProgramStart, 256); err != nil {
        t.Error(err)
    } else if checkMemory(cpu, ProgramStart, 256) != nil {
        t.Error(err)
    }
    if err := cpu.SetMemory(DefaultMemorySize, 256); err == nil {
        t.Error("Memory access beyond registered size.")
    }
    if _, err := cpu.GetMemory(DefaultMemorySize); err == nil {
        t.Error("Memory access beyond registered size.")
    }
}

func testInstructionHelper(t *testing.T, op Opcode, src uint16, tgt uint16, memloc uint16, exp uint16) {
    inst := NewInstruction(op, src, tgt, memloc)
    if inst.Encode() != exp {
        t.Error(fmt.Sprintf("Expected (%d) = %b, was %b.", op, exp, inst.Encode()))
    }
    instd := DecodeInstruction(inst.Encode())
    if instd.opcode != op {
        t.Error(fmt.Sprintf("Invalid decoding, expected %d, got %d.", op, instd.opcode))
    }
    if instd.regsrc != src || instd.regtgt != tgt || instd.memloc != memloc {
        t.Error(fmt.Sprintf("Parameters did not decode: %d/%d, %d/%d, %d/%d.",
            src, instd.regsrc, tgt, instd.regtgt, memloc, instd.memloc))
    }
}

func TestInstructions(t *testing.T) {
    // Test that converting an instruction to an integer produces the correct bit pattern.
    testInstructionHelper(t, OpHalt,   0, 0, 0, 0b0000000000000000)
    testInstructionHelper(t, OpNoop,   0, 0, 0, 0b0001000000000000)
    testInstructionHelper(t, OpSetReg, 0, 1, 0, 0b1000100000000001)
}

func TestBcpuNoop(t *testing.T) {
    var cpu *Bcpu = NewBcpu()
    noop := NewInstruction(OpNoop,0,0,0).Encode()
    cpu.SetMemory(ProgramStart, noop)
    cpu.SetMemory(ProgramStart, noop)
    if err := cpu.Run(); err != nil {
        t.Error(fmt.Sprintf("Execution error: %s.", err))
    }
    exppc := ProgramStart + 2
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

func checkMemory(cpu *Bcpu, location uint16, expval uint16) error {
    if val, err := cpu.GetMemory(location); err != nil {
        return err
    } else {
        if val != expval {
            return fmt.Errorf("Location %d should have a value of %d, has a value of %d.",
                ProgramStart, expval, val)
        }
    }
    return nil
}

func testRegisterSet(t *testing.T, cpu *Bcpu, reg uint16, expval uint16) {
    sr := NewInstruction(OpSetReg,0,reg,0).Encode()
    cpu.SetMemory(ProgramStart, sr)
    if err := checkMemory(cpu, ProgramStart, sr); err != nil {
        t.Error(err)
    }
    cpu.SetMemory(ProgramStart+1, expval)
    if err := checkMemory(cpu, ProgramStart+1, expval); err != nil {
        t.Error(err)
    }
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
