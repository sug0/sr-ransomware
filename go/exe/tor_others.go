// +build !windows

package exe

import "os/exec"

func (t *Tor) Start() error {
    if t.cmd != nil {
        return ErrAlreadyRunning
    }
    if t.config != "" {
        t.cmd = exec.Command("tor", "-f", t.config)
    } else {
        t.cmd = exec.Command("tor")
    }
    return t.bootstrap()
}
