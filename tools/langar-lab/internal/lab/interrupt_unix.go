//go:build !windows

package lab

import "os"

func interruptProcess(p *os.Process) error { return p.Signal(os.Interrupt) }
