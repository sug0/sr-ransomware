package fs

import (
    "os"
    "io"
    "archive/zip"
    "path/filepath"

    "github.com/sug0/sr-ransomware/go/errors"
)

func AllDrives() (drives []string) {
    for _, letter := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
        drive := string(letter) + `:\`
        if _, err := os.Stat(drive); err == nil {
            drives = append(drives, drive)
        }
    }
    return
}

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

func Unzip(to, from string) error {
    return errors.WrapIfNotNil(pkg, "failed to unzip file", unzip(to, from))
}

// https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func unzip(dest, src string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer r.Close()

    err = os.MkdirAll(dest, 0755)
    if err != nil && !os.IsExist(err) {
        return err
    }

    extractAndWriteFile := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer rc.Close()

        path := filepath.Join(dest, f.Name)

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer f.Close()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return err
        }
    }

    return nil
}
