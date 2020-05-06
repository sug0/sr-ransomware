package util

import (
    "fmt"
    "crypto/rsa"
    "crypto/x509"
    //"crypto/rand"
    //"encoding/pem"
)

const pkg = "github.com/sug0/sr-ransomware/go/crypto/util"

var (
    errNotRSA = fmt.Errorf("%s: key is not RSA", pkg)
)

// Parses a DER encoded *rsa.PrivateKey and returns it.
func ImportDERSecretKeyRSA(data []byte) (*rsa.PrivateKey, error) {
    return x509.ParsePKCS1PrivateKey(data)
}

// Parses a DER encoded *rsa.PublicKey and returns it.
func ImportDERPublicKeyRSA(data []byte) (*rsa.PublicKey, error) {
    key, err := x509.ParsePKIXPublicKey(data)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to parse RSA Public Key: %w", pkg, err)
    }
    if k, ok := key.(*rsa.PublicKey); ok {
        return k, nil
    }
    return nil, errNotRSA
}
