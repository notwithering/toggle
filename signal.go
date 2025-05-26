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
	var value string
	err := ctx.Scan.PopValueInto("signal", &value)
	if err != nil {
		return err
	}

	var errUnknownSignal = fmt.Errorf("unknown signal %s", value)

	n, err := strconv.ParseInt(value, 10, 0)
	if err == nil {
		if n <= 0 || n > 31 {
			return errUnknownSignal
		}
		*s = KongSignal(n)
		return nil
	}

	// Formatted from builtin/syscall/zerrors_linux_amd64.go "Signals" const block.
	switch strings.TrimPrefix(strings.ToUpper(value), "SIG") {
	case "ABRT":
		*s = 0x6
	case "ALRM":
		*s = 0xe
	case "BUS":
		*s = 0x7
	case "CHLD":
		*s = 0x11
	case "CLD":
		*s = 0x11
	case "CONT":
		*s = 0x12
	case "FPE":
		*s = 0x8
	case "HUP":
		*s = 0x1
	case "ILL":
		*s = 0x4
	case "INT":
		*s = 0x2
	case "IO":
		*s = 0x1d
	case "IOT":
		*s = 0x6
	case "KILL":
		*s = 0x9
	case "PIPE":
		*s = 0xd
	case "POLL":
		*s = 0x1d
	case "PROF":
		*s = 0x1b
	case "PWR":
		*s = 0x1e
	case "QUIT":
		*s = 0x3
	case "SEGV":
		*s = 0xb
	case "STKFLT":
		*s = 0x10
	case "STOP":
		*s = 0x13
	case "SYS":
		*s = 0x1f
	case "TERM":
		*s = 0xf
	case "TRAP":
		*s = 0x5
	case "TSTP":
		*s = 0x14
	case "TTIN":
		*s = 0x15
	case "TTOU":
		*s = 0x16
	case "UNUSED":
		*s = 0x1f
	case "URG":
		*s = 0x17
	case "USR1":
		*s = 0xa
	case "USR2":
		*s = 0xc
	case "VTALRM":
		*s = 0x1a
	case "WINCH":
		*s = 0x1c
	case "XCPU":
		*s = 0x18
	case "XFSZ":
		*s = 0x19
	default:
		return errUnknownSignal
	}

	return nil
}
