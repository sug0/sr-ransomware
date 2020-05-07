package fs

import (
    "os"
    "io"

    "github.com/sug0/sr-ransomware/go/errors"
)

func Move(to, from string) error {
    err := move(to, from)
    if err != nil {
        return err
    }
    return errors.WrapIfNotNil(pkg, "failed to remove file", os.Remove(from))
}

func move(to, from string) error {
    fto, err := os.Create(to)
    if err != nil {
        return errors.Wrap(pkg, "failed to create file", err)
    }
    defer fto.Close()

    ffrom, err := os.Open(from)
    if err != nil {
        return errors.Wrap(pkg, "failed to open file", err)
    }
    defer ffrom.Close()

    _, err = io.Copy(fto, ffrom)
    return errors.WrapIfNotNil(pkg, "failed to copy file to new path", err)
}
