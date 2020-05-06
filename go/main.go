package main

import (
    "fmt"

    "github.com/sug0/sr-ransomware/go/crypto/util"
)

func main() {
    key, err := util.ParseDERPublicKeyRSA([]byte{1,2,3,4})
    if err != nil {
        panic(err)
    }
    fmt.Printf("%#v\n", key)
}
