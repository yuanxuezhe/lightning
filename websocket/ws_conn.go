package network

import (
	"errors"
	"github.com/gorilla/websocket"
	"net"
	"sync"
)

type WebsocketConnSet map[*websocket.Conn]struct{}

type WSConn struct {
	sync.Mutex
	conn *websocket.Conn
	//writeChan chan []byte
	maxMsgLen uint32
	closeFlag bool
}

func newWSConn(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *WSConn {
	wsConn := new(WSConn)
	wsConn.conn = conn
	//wsConn.writeChan = make(chan []byte, pendingWriteNum)
	//wsConn.maxMsgLen = maxMsgLen
	//
	//go func() {
	//	for b := range wsConn.writeChan {
	//		if b == nil {
	//			break
	//		}
	//
	//		err := conn.WriteMessage(websocket.BinaryMessage, b)
	//		if err != nil {
	//			break
	//		}
	//	}
	//
	//	conn.Close()
	//	wsConn.Lock()
	//	wsConn.closeFlag = true
	//	wsConn.Unlock()
	//}()

	return wsConn
}

func (wsConn *WSConn) doDestroy() {
	wsConn.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	wsConn.conn.Close()

	if !wsConn.closeFlag {
		//close(wsConn.writeChan)
		wsConn.closeFlag = true
	}
}

func (wsConn *WSConn) Destroy() {
	wsConn.Lock()
	defer wsConn.Unlock()

	wsConn.doDestroy()
}

func (wsConn *WSConn) Close() error {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return nil
	}

	//wsConn.doWrite(nil)
	wsConn.closeFlag = true
	return nil
}

//
//func (wsConn *WSConn) doWrite(b []byte) {
//	if len(wsConn.writeChan) == cap(wsConn.writeChan) {
//		log.Println("close conn: channel full")
//		wsConn.doDestroy()
//		return
//	}
//
//	wsConn.writeChan <- b
//}

func (wsConn *WSConn) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

func (wsConn *WSConn) RemoteAddr() net.Addr {
	return wsConn.conn.RemoteAddr()
}

// goroutine not safe
func (wsConn *WSConn) ReadMsg() ([]byte, error) {
	_, b, err := wsConn.conn.ReadMessage()
	return b, err
}

// args must not be modified by the others goroutines
func (wsConn *WSConn) WriteMsg(arg []byte) error {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return nil
	}

	var msgLen uint32
	msgLen = uint32(len(arg))

	wsConn.maxMsgLen = 1<<(32) - 1

	if msgLen > wsConn.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < 1 {
		return errors.New("message too short")
	}
	// don't copy
	err := wsConn.conn.WriteMessage(websocket.BinaryMessage, arg)

	if err != nil {
		return err
	}
	return nil
}
