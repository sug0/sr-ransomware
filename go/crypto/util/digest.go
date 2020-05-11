package util

import (
    "crypto/sha1"
    "encoding/hex"
)

func Sha1Digest(data []byte) string {
    hexdigest := make([]byte, 2*sha1.Size)
    digest := sha1.Sum(data)
    hex.Encode(hexdigest, digest[:])
    return string(hexdigest)
}
