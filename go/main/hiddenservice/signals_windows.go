// +build windows

package main

import (
    "os"
    "os/signal"
)

func signalListener() <-chan os.Signal {
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, os.Kill, os.Interrupt)
    return ch
}

