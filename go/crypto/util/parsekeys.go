package util

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"

    "github.com/sug0/sr-ransomware/go/errors"
)

// Parses a DER encoded *rsa.PrivateKey and returns it.
func ParseDERSecretKeyRSA(data []byte) (*rsa.PrivateKey, error) {
    return x509.ParsePKCS1PrivateKey(data)
}

// Parses a DER encoded *rsa.PublicKey and returns it.
func ParseDERPublicKeyRSA(data []byte) (*rsa.PublicKey, error) {
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
func ParsePEMSecretKeyRSA(data []byte) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode(data)
    if block == nil {
        return nil, errNotPEM
    }
    if block.Type != "RSA PRIVATE KEY" {
        return nil, errNotRSA
    }
    return ParseDERSecretKeyRSA(block.Bytes)
}

// Parses a PEM encoded *rsa.PublicKey and returns it.
func ParsePEMPublicKeyRSA(data []byte) (*rsa.PublicKey, error) {
    block, _ := pem.Decode(data)
    if block == nil {
        return nil, errNotPEM
    }
    if block.Type != "PUBLIC KEY" {
        return nil, errNotPUB
    }
    return ParseDERPublicKeyRSA(block.Bytes)
}
