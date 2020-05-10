package ratelimit

import (
    "net"
    "net/http"
    "net/url"
    "context"
    "time"

    "github.com/sug0/sr-ransomware/go/errors"
)

// Instantiates a new http.Client with a ratelimit.Conn socket.
func NewHTTPClient(timeout, sleep time.Duration, proxyViaTor bool) http.Client {
    return http.Client{
        Transport: &http.Transport{
            DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
                var d net.Dialer
                conn, err := d.DialContext(ctx, network, addr)
                if err != nil {
                    return nil, errors.Wrap(pkg, "failed to dial address", err)
                }
                return NewConn(conn, sleep), nil
            },
            Proxy: getTorProxy(proxyViaTor),
        },
        Timeout: timeout,
    }
}

func getTorProxy(proxyViaTor bool) func(*http.Request) (*url.URL, error) {
    if proxyViaTor {
        const torSocksProxy = "socks5://127.0.0.1:9050"
        proxyURL, _ := url.Parse(torSocksProxy)
        return http.ProxyURL(proxyURL)
    }
    return nil
}
