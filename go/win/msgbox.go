// +build windows

package win

import (
    "unsafe"
    "syscall"
)

const (
    MB_OK                = 0x00000000
    MB_OKCANCEL          = 0x00000001
    MB_ABORTRETRYIGNORE  = 0x00000002
    MB_YESNOCANCEL       = 0x00000003
    MB_YESNO             = 0x00000004
    MB_RETRYCANCEL       = 0x00000005
    MB_CANCELTRYCONTINUE = 0x00000006
    MB_ICONHAND          = 0x00000010
    MB_ICONQUESTION      = 0x00000020
    MB_ICONEXCLAMATION   = 0x00000030
    MB_ICONASTERISK      = 0x00000040
    MB_USERICON          = 0x00000080
    MB_ICONWARNING       = MB_ICONEXCLAMATION
    MB_ICONERROR         = MB_ICONHAND
    MB_ICONINFORMATION   = MB_ICONASTERISK
    MB_ICONSTOP          = MB_ICONHAND

    MB_DEFBUTTON1 = 0x00000000
    MB_DEFBUTTON2 = 0x00000100
    MB_DEFBUTTON3 = 0x00000200
    MB_DEFBUTTON4 = 0x00000300
)

const (
    IDOK       = 1
    IDCANCEL   = 2
    IDABORT    = 3
    IDRETRY    = 4
    IDIGNORE   = 5
    IDYES      = 6
    IDNO       = 7
    IDCLOSE    = 8
    IDHELP     = 9
    IDTRYAGAIN = 10
    IDCONTINUE = 11
    IDTIMEOUT  = 32000
)

var (
    user32dll *syscall.LazyDLL
    msgb      *syscall.LazyProc
)

func init() {
    user32dll = syscall.NewLazyDLL("user32.dll")
    msgb = user32dll.NewProc("MessageBoxW")
}

func MessageBox(title, message string, flags int) int {
    t := uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title)))
    m := uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(message)))
    id, _, _ := msgb.Call(0, m, t, uintptr(flags))
    return int(id)
}
