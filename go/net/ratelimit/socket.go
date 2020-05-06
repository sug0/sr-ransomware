package ratelimit

import (
    "net"
    "time"
)

// Represents a rate limited net.Conn socket.
type Conn struct {
    sleep time.Duration
    conn  net.Conn
}

// Returns a new rate limited socket. Before any `Read` or `Write`
// operation, the current go routine sleeps for `sleep` time duration.
//
// Furthermore, the I/O buffers are limited to 4 KiB in size; this way,
// if `sleep = 1*time.Second`, then the rw speed is limited to roughly 4 KiB/s.
func NewConn(conn net.Conn, sleep time.Duration) net.Conn {
    return &Conn{
        sleep: sleep,
        conn: conn,
    }
}

// Implements the net.Conn interface.
func (c *Conn) Read(p []byte) (int, error) {
    if len(p) > 4096 {
        p = p[:4096]
    }
    time.Sleep(c.sleep)
    return c.conn.Read(p)
}

// Implements the net.Conn interface.
func (c *Conn) Write(p []byte) (int, error) {
    if len(p) > 4096 {
        p = p[:4096]
    }
    time.Sleep(c.sleep)
    return c.conn.Write(p)
}

// Implements the net.Conn interface.
func (c *Conn) Close() error {
    return c.conn.Close()
}

// Implements the net.Conn interface.
func (c *Conn) LocalAddr() net.Addr {
    return c.conn.LocalAddr()
}

// Implements the net.Conn interface.
func (c *Conn) RemoteAddr() net.Addr {
    return c.conn.RemoteAddr()
}

// Implements the net.Conn interface.
func (c *Conn) SetDeadline(t time.Time) error {
    return c.conn.SetDeadline(t)
}

// Implements the net.Conn interface.
func (c *Conn) SetReadDeadline(t time.Time) error {
    return c.conn.SetReadDeadline(t)
}

// Implements the net.Conn interface.
func (c *Conn) SetWriteDeadline(t time.Time) error {
    return c.conn.SetWriteDeadline(t)
}
