package win

import "github.com/sug0/sr-ransomware/go/errors"

const pkg = "github.com/sug0/sr-ransomware/go/win"

var (
    errNotLoggedIn = errors.New(pkg, "no user logged in")
)
