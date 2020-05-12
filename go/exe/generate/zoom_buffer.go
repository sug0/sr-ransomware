package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "bufio"
    "net/http"
    "io/ioutil"
    "path/filepath"
)

const workdir = "tmp"
var zoomPath = filepath.Join(workdir, "ZoomInstaller.exe")

func main() {
    log.Println("> Generating github.com/sug0/sr-ransomware/go/exe/zoom_buffer.go")
    if _, err := os.Stat("zoom_buffer.go"); err == nil {
        return
    }
    f, err := os.Create("zoom_buffer.go")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    zoom, err := zoomBytes()
    if err != nil {
        log.Fatal(err)
    }
    w := bufio.NewWriter(f)
    fmt.Fprintf(w, "// +build zoom\n\npackage exe;var zoomEXE=[]byte{")
    for i := 0; i < len(zoom); i++ {
        fmt.Fprintf(w, "%d,", zoom[i])
    }
    fmt.Fprintf(w, "}")
    w.Flush()
}

func zoomBytes() ([]byte, error) {
    if err := os.Mkdir(workdir, os.ModePerm); err != nil && !os.IsExist(err) {
        return nil, err
    }
    if _, err := os.Stat(zoomPath); err != nil {
        err = downloadZoom()
        if err != nil {
            return nil, err
        }
    }
    f, err := os.Open(zoomPath)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return ioutil.ReadAll(f)
}

func downloadZoom() error {
    f, err := os.Create(zoomPath)
    if err != nil {
        return err
    }
    defer f.Close()
    rsp, err := http.Get("https://zoom.us/client/latest/ZoomInstaller.exe")
    if err != nil {
        return err
    }
    defer rsp.Body.Close()
    _, err = io.Copy(f, rsp.Body)
    return err
}
