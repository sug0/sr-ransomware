package util

import "github.com/sug0/sr-ransomware/go/errors"

const pkg = "github.com/sug0/sr-ransomware/go/crypto/util"

var (
    errNotRSA = errors.New(pkg, "key is not RSA")
    errNotPEM = errors.New(pkg, "invalid PEM block")
    errNotPUB = errors.New(pkg, "key is not public")
)
