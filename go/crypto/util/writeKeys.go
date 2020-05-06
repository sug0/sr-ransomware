package util

import (
    "crypto/rsa"
    "crypto/x509"

    "github.com/sug0/sr-ransomware/go/errors"
)

// Marshals a DER encoded *rsa.PrivateKey and returns it.
func ExportDERSecretKeyRSA(key *rsa.PrivateKey) ([]byte, error) {
    data := x509.MarshalPKCS1PrivateKey(key)
    if data == nil {
        return nil, errNotRSA
    }
    return data, nil
}

// Marshals a DER encoded *rsa.PublicKey and returns it.
func ExportDERPublicKeyRSA(key *rsa.PublicKey) ([]byte, error) {
    data, err := x509.MarshalPKIXPublicKey(key)
    return data, errors.WrapIfNotNil(pkg, "invalid RSA key", err)
}
