package main

import (
    "os"
    "log"
    "fmt"
    "time"
    "io/ioutil"
    "path/filepath"

    "github.com/sug0/sr-ransomware/go/crypto/util"
)

func main() {
    t := time.Now()
    log.Println("> Generating github.com/sug0/sr-ransomware/go/crypto/scheme/attacker/public.go")
    defer log.Printf("< Completed in %s\n", time.Since(t))

    pem, err := ioutil.ReadFile(filepath.Join("..", "..", "..", "..", "res", "attacker", "key.pub"))
    if err != nil {
        log.Fatal(err)
    }
    pk, err := util.ImportPEMPublicKeyRSA(pem)
    if err != nil {
        log.Fatal(err)
    }
    der, err := util.ExportDERPublicKeyRSA(pk)
    if err != nil {
        log.Fatal(err)
    }
    f, err := os.Create("public.go")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    fmt.Fprintf(f, "package attacker;")
    fmt.Fprintf(f, "var oraclePublicKey=[]byte{")
    for i := 0; i < len(der); i++ {
        fmt.Fprintf(f, "%d,", der[i])
    }
    fmt.Fprintf(f, "}")
}
