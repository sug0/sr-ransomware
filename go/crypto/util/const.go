package util

import "github.com/sug0/sr-ransomware/go/errors"

const pkg = "github.com/sug0/sr-ransomware/go/crypto/util"

var (
    errNotRSA = errors.New(pkg, "invalid RSA key")
    errNotPEM = errors.New(pkg, "invalid PEM block")
    errNotPUB = errors.New(pkg, "invalid public key")
)
