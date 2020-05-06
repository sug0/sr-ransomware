package ratelimit

import (
    "net"
    "net/http"
    "net/url"
    "context"
    "time"
)

// Instantiates a new http.Client with a ratelimit.Conn socket.
func NewHTTPClient(sleep time.Duration, proxyViaTOR bool) http.Client {
    var proxy func(*http.Request) (*url.URL, error)
    if proxyViaTOR {
        const torSocksProxy = "socks5://127.0.0.1:9050"
        proxyURL, _ := url.Parse(torSocksProxy)
        proxy = http.ProxyURL(proxyURL)
    }
    return http.Client{
        Transport: &http.Transport{
            DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
                var d net.Dialer
                conn, err := d.DialContext(ctx, network, addr)
                if err != nil {
                    return nil, err
                }
                return NewConn(conn, sleep), nil
            },
            Proxy: proxy,
        },
    }
}
