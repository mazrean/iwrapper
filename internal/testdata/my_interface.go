package testdata

import (
	"bufio"
	"net"
)

type Hijacker interface {
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}
