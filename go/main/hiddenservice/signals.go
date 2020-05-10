// +build !windows

package main

import (
    "os"
    "syscall"
    "os/signal"
)

func signalListener() <-chan os.Signal {
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
    return ch
}
