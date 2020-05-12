// +build windows

package win

import "syscall"

var (
    user32DLL   *syscall.LazyDLL
    messageBoxW *syscall.LazyProc
)

func init() {
    user32DLL = syscall.NewLazyDLL("user32.dll")
    messageBoxW = user32DLL.NewProc("MessageBoxW")
}
