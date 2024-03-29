package network

import (
	"log"
	"net"
	"sync"
	"time"

	conn "github.com/yuanxuezhe/lightning/Conn"
)

type TCPClient struct {
	sync.Mutex
	Addr            string
	ConnNum         int
	ConnectInterval time.Duration
	PendingWriteNum int
	AutoReconnect   bool
	conns           ConnSet
	wg              sync.WaitGroup
	closeFlag       bool

	// msg parser
	LenMsgLen    int
	MinMsgLen    uint32
	MaxMsgLen    uint32
	LittleEndian bool
	//msgParser    *MsgParser
}

//func (client *TCPClient) Start() {
//	client.init()
//
//	for i := 0; i < client.ConnNum; i++ {
//		client.wg.Add(1)
//		go client.connect()
//	}
//}

//func (client *TCPClient) init() {
//	client.Lock()
//	defer client.Unlock()
//
//	if client.ConnNum <= 0 {
//		client.ConnNum = 1
//		log.Printf("invalid ConnNum, reset to %v", client.ConnNum)
//	}
//	if client.ConnectInterval <= 0 {
//		client.ConnectInterval = 3 * time.Second
//		log.Printf("invalid ConnectInterval, reset to %v", client.ConnectInterval)
//	}
//	if client.PendingWriteNum <= 0 {
//		client.PendingWriteNum = 100
//		log.Printf("invalid PendingWriteNum, reset to %v", client.PendingWriteNum)
//	}
//
//	if client.conns != nil {
//		log.Fatal("client is running")
//	}
//
//	client.conns = make(ConnSet)
//	client.closeFlag = false
//
//	// msg parser
//	//msgParser := NewMsgParser()
//	//msgParser.SetMsgLen(client.LenMsgLen, client.MinMsgLen, client.MaxMsgLen)
//	//msgParser.SetByteOrder(client.LittleEndian)
//	//client.msgParser = msgParser
//}

func (client *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil || client.closeFlag {
			return conn
		}

		log.Printf("connect to %v error: %v", client.Addr, err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *TCPClient) Connect() conn.CommConn {
	defer client.Close()
	//reconnect:
	conn := client.dial()
	if conn == nil {
		return nil
	}

	client.Lock()
	if client.closeFlag {

		client.Unlock()
		conn.Close()
		return nil
	}

	client.Unlock()

	tcpConn := newTCPConn(conn, client.PendingWriteNum)

	return tcpConn
	// cleanup
	//tcpConn.Close()
	//client.Lock()
	//delete(client.conns, conn)
	//client.Unlock()
	//agent.OnClose()

	//if client.AutoReconnect {
	//	time.Sleep(client.ConnectInterval)
	//	goto reconnect
	//}
}

func (client *TCPClient) Close() {
	client.Lock()
	client.closeFlag = true
	//for conn := range client.conns {
	//	conn.Close()
	//}
	//client.conns = nil
	client.Unlock()

	//client.wg.Wait()
}
