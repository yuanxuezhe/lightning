package network

import (
	conn "gitee.com/yuanxuezhe/ynet/Conn"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type WSClient struct {
	sync.Mutex
	Addr             string
	ConnNum          int
	ConnectInterval  time.Duration
	PendingWriteNum  int
	MaxMsgLen        uint32
	HandshakeTimeout time.Duration
	AutoReconnect    bool
	dialer           websocket.Dialer
	conns            WebsocketConnSet
	//wg               sync.WaitGroup
	closeFlag bool
}

func (client *WSClient) dial() *websocket.Conn {
	for {
		conn, _, err := client.dialer.Dial(client.Addr, nil)
		if err == nil || client.closeFlag {
			return conn
		}
		log.Printf("connect to %v error: %v", client.Addr, err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *WSClient) Connect() conn.CommConn {
	defer client.Close()

	//reconnect:
	conn := client.dial()
	if conn == nil {
		//if client.AutoReconnect {
		//	time.Sleep(client.ConnectInterval)
		//	goto reconnect
		//}
		return nil
	}

	conn.SetReadLimit(int64(client.MaxMsgLen))

	client.Lock()
	if client.closeFlag {
		client.Unlock()
		conn.Close()
		return nil
	}

	client.Unlock()

	wsConn := newWSConn(conn, client.PendingWriteNum, client.MaxMsgLen)

	return wsConn
}

func (client *WSClient) Close() {
	client.Lock()
	client.closeFlag = true
	for conn := range client.conns {
		conn.Close()
	}
	client.conns = nil
	client.Unlock()

	//client.wg.Wait()
}
