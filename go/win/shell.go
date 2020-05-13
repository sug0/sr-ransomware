// +build windows

package win

import (
    "os"
    "syscall"

    "golang.org/x/sys/windows"
    "github.com/sug0/sr-ransomware/go/errors"
)

const (
    SW_HIDE            = 0
    SW_NORMAL          = 1
    SW_SHOWNORMAL      = 1
    SW_SHOWMINIMIZED   = 2
    SW_MAXIMIZE        = 3
    SW_SHOWMAXIMIZED   = 3
    SW_SHOWNOACTIVATE  = 4
    SW_SHOW            = 5
    SW_MINIMIZE        = 6
    SW_SHOWMINNOACTIVE = 7
    SW_SHOWNA          = 8
    SW_RESTORE         = 9
    SW_SHOWDEFAULT     = 10
    SW_FORCEMINIMIZE   = 11
)

var (
    shell32 = syscall.NewLazyDLL("Shell32.dll")
    isAdmin = shell32.NewProc("IsUserAnAdmin")
)

func ShellExecute(lpOperation, lpFile string, nShowCmd int) error {
    err := windows.ShellExecute(
        0,
        syscall.StringToUTF16Ptr(lpOperation),
        syscall.StringToUTF16Ptr(lpFile),
        nil,
        nil,
        int32(nShowCmd))
    return errors.WrapIfNotNil(pkg, "shell exec failed", err)
}

func IsUserAnAdmin() bool {
    isAdmin, _, _ := syscall.Syscall(isAdmin.Addr(), 0, 0, 0, 0)
    return isAdmin == 1
}

func RunAsAdmin() {
    if !IsUserAnAdmin() {
        ShellExecute("runas", `"`+os.Args[0]+`"`, SW_SHOW)
        os.Exit(0)
    }
}
