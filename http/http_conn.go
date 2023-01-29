package network

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

type Addrs struct {
	network string
	ip      string
}

func (addr *Addrs) Network() string {
	return addr.network
}

func (addr *Addrs) String() string {
	return addr.ip
}

type HttpConn struct {
	w         http.ResponseWriter
	r         *http.Request
	writeChan chan []byte
	maxMsgLen uint32
}

func newHttpConn(w http.ResponseWriter, r *http.Request, pendingWriteNum int, maxMsgLen uint32) *HttpConn {
	httpConn := new(HttpConn)
	httpConn.w = w
	httpConn.r = r
	httpConn.writeChan = make(chan []byte, pendingWriteNum)
	httpConn.maxMsgLen = maxMsgLen

	go func() {
		for b := range httpConn.writeChan {
			if b == nil {
				break
			}

			//err := w.WriteMessage(websocket.BinaryMessage, b)
			_, err := w.Write(b)

			if err != nil {
				break
			}
		}
	}()

	return httpConn
}

func (conn *HttpConn) doDestroy() {
}

func (conn *HttpConn) Destroy() {
}

func (conn *HttpConn) Close() error {
	return nil
}

func (conn *HttpConn) LocalAddr() net.Addr {
	return &Addrs{
		"tcp",
		conn.r.Host,
	}
}

func (conn *HttpConn) RemoteAddr() net.Addr {
	return &Addrs{
		"tcp",
		conn.r.RemoteAddr,
	}
}

// goroutine not safe
func (conn *HttpConn) ReadMsg() ([]byte, error) {
	var err error
	err = nil

	query := conn.r.URL.Query()

	// 第一种方式
	// id := query["id"][0]

	// 第二种方式
	param := query.Get("param")

	if len(param) < 1 {
		err = errors.New("param is NULL")
	}

	return []byte(param), err
}

// args must not be modified by the others goroutines
func (conn *HttpConn) WriteMsg(arg []byte) error {
	var msgLen uint32
	msgLen = uint32(len(arg))

	conn.maxMsgLen = 1<<(32) - 1

	if msgLen > conn.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < 1 {
		return errors.New("message too short")
	}

	// don't copy
	_, err := conn.w.Write(arg)
	fmt.Println("WriteMsg(arg []byte):", string(arg))
	if err != nil {
		fmt.Println("err mssg  ", err.Error())
		return err
	}
	return nil
}
