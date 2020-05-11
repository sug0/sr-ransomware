package attacker

import "github.com/sug0/sr-ransomware/go/errors"

const pkg = "github.com/sug0/sr-ransomware/go/crypto/scheme/attacker"

const (
    rsaKeyBits  = 2048
    ransomValue = 0.29
)

var (
    ErrRansom = errors.New(pkg, "ransom not paid")
)
