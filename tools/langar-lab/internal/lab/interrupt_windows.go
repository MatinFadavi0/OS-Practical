package lab

import "os"

// Multipass is cancelled separately through the Linux worker on Windows.
func interruptProcess(p *os.Process) error { return p.Kill() }
