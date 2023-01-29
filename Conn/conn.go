package conn

import "net"

type CommConn interface {
	ReadMsg() ([]byte, error)
	WriteMsg([]byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() error
	Destroy()
}
