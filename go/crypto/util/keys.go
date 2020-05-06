package util

import (
    "errors"
    "crypto/rsa"
    "crypto/x509"
    "crypto/rand"
    "encoding/pem"
)

func ImportDERSecretKeyRSA(data []byte) (*rsa.PrivateKey, error) {
    return x509.ParsePKCS1PrivateKey(data)
}

func ImportDERPublicKeyRSA(data []byte) (*rsa.PublicKey, error) {
    key, err := x509.ParsePKIXPublicKey(data)
    if err != nil {
        return nil, err
    }
    if k, ok := key.(*rsa.PublicKey); ok {
        return k, nil
    }
    return nil, errors.New("key is not RSA")
}
