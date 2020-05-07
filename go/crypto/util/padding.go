package util

import "bytes"

func Pad(src []byte, blockSize int) []byte {
    // https://gist.github.com/stupidbodo/601b68bfef3449d1b8d9
    padding := blockSize - len(src)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(src, padtext...)
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
