package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "os/exec"
    "net/http"
    "io/ioutil"
    "path/filepath"
)

func main() {
    if _, err := os.Stat("tor_buffer.go"); err == nil {
        return
    }
    f, err := os.Create("tor_buffer.go")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    tor, err := torBytes()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Fprintf(f, "package exe;")
    fmt.Fprintf(f, "var torZIP=[]byte{")
    for i := 0; i < len(tor); i++ {
        fmt.Fprintf(f, "%d,", tor[i])
    }
    fmt.Fprintf(f, "}")
}

func torBytes() ([]byte, error) {
    if _, err := os.Stat("TorInstaller.exe"); err != nil {
        err = downloadTor()
        if err != nil {
            return nil, err
        }
    }
    if _, err := os.Stat("Tor.zip"); err != nil {
        err = packTor()
        if err != nil {
            return nil, err
        }
    }
    f, err := os.Open("TorInstaller.exe")
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return ioutil.ReadAll(f)
}

func packTor() error {
    err := exec.Command("7z", "x", "TorInstaller.exe").Run()
    if err != nil {
        return err
    }
    return exec.Command(
        "7z",
        "-tzip", "-m0=lzma", "-mx=9",
        "a", "Tor.zip",
        filepath.Join("Browser", "TorBrowser", "Tor"),
    )
}

func downloadTor() error {
    f, err := os.Create("TorInstaller.exe")
    if err != nil {
        return err
    }
    defer f.Close()
    rsp, err := http.Get("https://www.torproject.org/dist/torbrowser/9.0.10/torbrowser-install-9.0.10_en-US.exe")
    if err != nil {
        return err
    }
    defer rsp.Body.Close()
    _, err = io.Copy(f, rsp.Body)
    return err
}

