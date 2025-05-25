package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"github.com/alecthomas/kong"
)

type KongSignal syscall.Signal

func (s *KongSignal) Decode(ctx *kong.DecodeContext) error {
	var raw string
	if err := ctx.Scan.PopValueInto("signal", &raw); err != nil {
		return err
	}

	var errUnknownSignal = fmt.Errorf("unknown signal \"%s\"", raw)

	n, err := strconv.Atoi(raw)
	if err == nil {
		if *s < 0 || *s > 31 {
			return errUnknownSignal
		}
		*s = KongSignal(n)
		return nil
	}

	switch strings.ToUpper(raw) {
	case "ABRT":
		*s = KongSignal(0x6)
	case "ALRM":
		*s = KongSignal(0xe)
	case "BUS":
		*s = KongSignal(0x7)
	case "CHLD":
		*s = KongSignal(0x11)
	case "CLD":
		*s = KongSignal(0x11)
	case "CONT":
		*s = KongSignal(0x12)
	case "FPE":
		*s = KongSignal(0x8)
	case "HUP":
		*s = KongSignal(0x1)
	case "ILL":
		*s = KongSignal(0x4)
	case "INT":
		*s = KongSignal(0x2)
	case "IO":
		*s = KongSignal(0x1d)
	case "IOT":
		*s = KongSignal(0x6)
	case "KILL":
		*s = KongSignal(0x9)
	case "PIPE":
		*s = KongSignal(0xd)
	case "POLL":
		*s = KongSignal(0x1d)
	case "PROF":
		*s = KongSignal(0x1b)
	case "PWR":
		*s = KongSignal(0x1e)
	case "QUIT":
		*s = KongSignal(0x3)
	case "SEGV":
		*s = KongSignal(0xb)
	case "STKFLT":
		*s = KongSignal(0x10)
	case "STOP":
		*s = KongSignal(0x13)
	case "SYS":
		*s = KongSignal(0x1f)
	case "TERM":
		*s = KongSignal(0xf)
	case "TRAP":
		*s = KongSignal(0x5)
	case "TSTP":
		*s = KongSignal(0x14)
	case "TTIN":
		*s = KongSignal(0x15)
	case "TTOU":
		*s = KongSignal(0x16)
	case "UNUSED":
		*s = KongSignal(0x1f)
	case "URG":
		*s = KongSignal(0x17)
	case "USR1":
		*s = KongSignal(0xa)
	case "USR2":
		*s = KongSignal(0xc)
	case "VTALRM":
		*s = KongSignal(0x1a)
	case "WINCH":
		*s = KongSignal(0x1c)
	case "XCPU":
		*s = KongSignal(0x18)
	case "XFSZ":
		*s = KongSignal(0x19)
	default:
		return errUnknownSignal
	}

	return nil
}
