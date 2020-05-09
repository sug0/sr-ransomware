package util

import (
    "io"
    "crypto/rsa"
    "crypto/rand"

    "github.com/sug0/sr-ransomware/go/errors"
)

func GenerateKeyRSA(bits int) (*rsa.PrivateKey, error) {
    key, err := rsa.GenerateKey(rand.Reader, bits)
    return key, errors.WrapIfNotNil(pkg, "failed to generate RSA key", err)
}

func GenerateAES() ([]byte, error) {
    buf := make([]byte, 16)
    _, err := io.ReadFull(rand.Reader, buf)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to read AES key", err)
    }
    return buf, nil
}

func GenerateIVandKeyAES() ([]byte, error) {
    buf := make([]byte, 32)
    _, err := io.ReadFull(rand.Reader, buf)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to read AES IV+key", err)
    }
    return buf, nil
}
