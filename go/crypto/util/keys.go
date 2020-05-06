package util

import (
    "crypto/rsa"
    "crypto/x509"
    //"crypto/rand"
    //"encoding/pem"

    "github.com/sug0/sr-ransomware/go/errors"
)

const pkg = "github.com/sug0/sr-ransomware/go/crypto/util"

var (
    errNotRSA = errors.New(pkg, "key is not RSA")
)

// Parses a DER encoded *rsa.PrivateKey and returns it.
func ImportDERSecretKeyRSA(data []byte) (*rsa.PrivateKey, error) {
    return x509.ParsePKCS1PrivateKey(data)
}

// Parses a DER encoded *rsa.PublicKey and returns it.
func ImportDERPublicKeyRSA(data []byte) (*rsa.PublicKey, error) {
    key, err := x509.ParsePKIXPublicKey(data)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to parse RSA Public Key", err)
    }
    if k, ok := key.(*rsa.PublicKey); ok {
        return k, nil
    }
    return nil, errNotRSA
}
