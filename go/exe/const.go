package exe

import "github.com/sug0/sr-ransomware/go/errors"

const pkg = "github.com/sug0/sr-ransomware/go/crypto/exe"

var (
    ErrAlreadyRunning = errors.New(pkg, "already running")
)
