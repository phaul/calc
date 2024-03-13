package vm

import (
	"fmt"
	"log"
	"slices"
	"strconv"

	"github.com/paulsonkoly/calc/memory"
	"github.com/paulsonkoly/calc/types/bytecode"
	"github.com/paulsonkoly/calc/types/value"
)

type Type struct {
	ip int              // instruction pointer
	m  *memory.Type     // variables
	CS *[]bytecode.Type // CS is the code segment
	DS *[]value.Type    // DS is the data segment
}

func New(m *memory.Type, cs *[]bytecode.Type, ds *[]value.Type) *Type {
	return &Type{m: m, CS: cs, DS: ds}
}

func (vm *Type) Run(retResult bool) value.Type {
	m := vm.m
	ds := vm.DS
  cs := vm.CS
  ip := vm.ip

	for ip < len(*cs) {
		instr := (*cs)[ip]

		// TODO allow tracing flag
		// fmt.Printf("%8d | %v\n", vm.ip, instr)

		opCode := instr.OpCode()

		switch opCode {
		case bytecode.ADD, bytecode.SUB, bytecode.MUL, bytecode.DIV:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			src1 := vm.fetch(instr.Src1(), instr.Src1Addr(), m, ds)
			op := [...]string{"+", "-", "*", "/"}[opCode-bytecode.ADD]

			val := src1.Arith(op, src0)

			m.Push(val)

		case bytecode.MOD:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			src1 := vm.fetch(instr.Src1(), instr.Src1Addr(), m, ds)

			val := src1.Mod(src0)

			m.Push(val)

		case bytecode.AND, bytecode.OR:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			src1 := vm.fetch(instr.Src1(), instr.Src1Addr(), m, ds)
			op := [...]string{"&", "|"}[opCode-bytecode.AND]

			val := src1.Logic(op, src0)

			m.Push(val)

		case bytecode.NOT:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			val := src0.Not()
			m.Push(val)

		case bytecode.LT, bytecode.GT, bytecode.LE, bytecode.GE:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			src1 := vm.fetch(instr.Src1(), instr.Src1Addr(), m, ds)
			op := [...]string{"<", ">", "<=", ">="}[opCode-bytecode.LT]

			val := src1.Relational(op, src0)

			m.Push(val)

		case bytecode.EQ, bytecode.NE:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			src1 := vm.fetch(instr.Src1(), instr.Src1Addr(), m, ds)
			op := [...]string{"==", "!="}[opCode-bytecode.EQ]

			val := src1.Eq(op, src0)

			m.Push(val)

		case bytecode.LEN:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)

			m.Push(src0.Len())

		case bytecode.IX1:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			src1 := vm.fetch(instr.Src1(), instr.Src1Addr(), m, ds)

			val := src1.Index(src0)

			m.Push(val)

		case bytecode.IX2:
      src2 := vm.fetch(instr.Src2(), instr.Src2Addr(), m, ds)
			src1 := vm.fetch(instr.Src1(), instr.Src1Addr(), m, ds)
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)

			val := src2.Index(src1, src0)

			m.Push(val)

		case bytecode.JMP:
			src0Imm := instr.Src0Addr()

			ip += src0Imm - 1

		case bytecode.JMPF:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			src1Imm := instr.Src1Addr()
			// TODO - how do we do type errors?
			if b, ok := src0.ToBool(); !b && ok {
				ip += src1Imm - 1
			}

		case bytecode.PUSH:
			src0 := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			m.Push(src0)

		case bytecode.POP:
			m.Pop()

		case bytecode.MOV:
			val := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			src1T := instr.Src1()

			switch src1T {
			case bytecode.ADDR_LCL:
				m.Set(instr.Src1Addr(), val)
			case bytecode.ADDR_GBL:
				name, ok := (*ds)[instr.Src1Addr()].ToString()
				if !ok {
					log.Panicf("unknown global\n %8d | %v\n", ip, instr)
				}
				m.SetGlobal(name, val)

			default:
				log.Panicf("unexpected dst in MOV\n %8d | %v\n", ip, instr)
			}

		case bytecode.ARR:
			val := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			ary := vm.fetch(instr.Src1(), instr.Src1Addr(), m, ds)

			slc, ok := ary.ToArray()
			if !ok {
				log.Panicf("cannot convert value to array\n %8d | %v\n", ip, instr)
			}

			slc = append(slc, val)

			val = value.NewArray(slc)
			m.Push(val)

		case bytecode.FUNC:
      localCnt := instr.Src2Addr()
      paramCnt := instr.Src1Addr()

			f := value.NewFunction(instr.Src0Addr(), slices.Clone(m.Top()), paramCnt, localCnt)
			m.Push(f)

		case bytecode.CALL:
			f := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			args := instr.Src1Addr()

			fVal, ok := f.ToFunction()

			if !ok {
				m.Push(value.TypeError)
				break
			}

			if fVal.ParamCnt != args {
				m.Push(value.ArgumentError)
				break
			}

			m.PushFrame(args, fVal.LocalCnt)
			m.PushClosure(fVal.Frame.(memory.Frame))
			m.Push(value.NewInt(ip))

			ip = fVal.Node - 1

		case bytecode.RET:
			lip, ok := m.IP().ToInt()
			if !ok {
				log.Panicf("can't pop instruction pointer\n %8d | %v\n", lip, instr)
			}

			val := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)

			m.PopFrame()
			m.PopClosure()

			m.Push(val)

			ip = lip

		case bytecode.WRITE:
			val := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			fmt.Println(val)
			m.Push(value.NoResultError)

		case bytecode.ATON:
			val := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)

			sv, ok := val.ToString()
			if !ok {
				m.Push(value.TypeError)
				break
			}

			if v, err := strconv.Atoi(string(sv)); err == nil {
				m.Push(value.NewInt(v))
				break
			}

			if v, err := strconv.ParseFloat(string(sv), 64); err == nil {
				m.Push(value.NewFloat(v))
				break
			}

			m.Push(value.ConversionError)

		case bytecode.TOA:
			val := vm.fetch(instr.Src0(), instr.Src0Addr(), m, ds)
			val = value.NewString(fmt.Sprint(val))
			m.Push(val)

		default:
			log.Panicf("unknown opcode: %v\n %8d | %v\n", opCode, ip, instr)
		}

		ip++
	}

  vm.ip = ip

	if retResult {
		return m.Pop()
	}

	return value.Type{}
}

func (vm Type) fetch(src uint64, addr int, m *memory.Type, ds *[]value.Type) value.Type {
	switch src {
	case bytecode.ADDR_STCK:
		return m.Pop()
	case bytecode.ADDR_DS:
		return (*ds)[addr]
	case bytecode.ADDR_CLS:
		return m.LookUpClosure(addr)
	case bytecode.ADDR_LCL:
		return m.LookUpLocal(addr)
	case bytecode.ADDR_GBL:
		name, ok := (*ds)[addr].ToString()
		if !ok {
			log.Panicf("unknown global\n %8d | %v\n", vm.ip, (*vm.CS)[vm.ip])
		}
		return m.LookUpGlobal(name)
	default:
		log.Panicf("unknown source\n %8d | %v\n", vm.ip, (*vm.CS)[vm.ip])
	}
	panic("unreachable code")
}
