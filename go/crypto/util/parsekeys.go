package util

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"

    "github.com/sug0/sr-ransomware/go/errors"
)

// Parses a DER encoded *rsa.PrivateKey and returns it.
func ImportDERSecretKeyRSA(data []byte) (*rsa.PrivateKey, error) {
    key, err := x509.ParsePKCS1PrivateKey(data)
    return key, errors.WrapIfNotNil(pkg, "import secret key failed", err)
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

// Parses a PEM encoded *rsa.PrivateKey and returns it.
func ImportPEMSecretKeyRSA(data []byte) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode(data)
    if block == nil {
        return nil, errNotPEM
    }
    if block.Type != "RSA PRIVATE KEY" {
        return nil, errNotRSA
    }
    return ImportDERSecretKeyRSA(block.Bytes)
}

// Parses a PEM encoded *rsa.PublicKey and returns it.
func ImportPEMPublicKeyRSA(data []byte) (*rsa.PublicKey, error) {
    block, _ := pem.Decode(data)
    if block == nil {
        return nil, errNotPEM
    }
    if block.Type != "PUBLIC KEY" {
        return nil, errNotPUB
    }
    return ImportDERPublicKeyRSA(block.Bytes)
}
