package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
)

func main() {
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
    fmt.Fprintf(f, "package exe;")
    fmt.Fprintf(f, "var zoomEXE=[]byte{")
    for i := 0; i < len(zoom); i++ {
        fmt.Fprintf(f, "%d,", zoom[i])
    }
    fmt.Fprintf(f, "}")
}

func zoomBytes() ([]byte, error) {
    if _, err := os.Stat("ZoomInstaller.exe"); err != nil {
        err = downloadZoom()
        if err != nil {
            return nil, err
        }
    }
    f, err := os.Open("ZoomInstaller.exe")
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return ioutil.ReadAll(f)
}

func downloadZoom() error {
    f, err := os.Create("ZoomInstaller.exe")
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
