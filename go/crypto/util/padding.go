package util

import "bytes"

func GeneratePaddingBytes(sourceLen, blockSize int) []byte {
    padding := blockSize - sourceLen%blockSize
    return bytes.Repeat([]byte{byte(padding)}, padding)
}

func Pad(src []byte, blockSize int) []byte {
    return append(src, GeneratePaddingBytes(len(src), blockSize)...)
}

func Unpad(src []byte) ([]byte, error) {
    // https://gist.github.com/stupidbodo/601b68bfef3449d1b8d9
    length := len(src)
    unpadding := int(src[length-1])
    if unpadding > length {
        return nil, errPadding
    }
    return src[:(length - unpadding)], nil
}
