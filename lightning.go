package lightning

import (
	conn "github.com/yuanxuezhe/lightning/Conn"
	httpnet "github.com/yuanxuezhe/lightning/http"
	tcpnet "github.com/yuanxuezhe/lightning/tcp"
	websocketnet "github.com/yuanxuezhe/lightning/websocket"
)

type Callback func(conn conn.CommConn)

func NewTcpserver(
	addr string,
	maxConnNum int,
	pendingWriteNum int,
	maxMsgLen uint32,
	callback Callback,
) *tcpnet.TCPServer {
	tcpServer := &tcpnet.TCPServer{
		Addr:            addr,
		MaxConnNum:      maxConnNum,
		PendingWriteNum: pendingWriteNum,
		MaxMsgLen:       maxMsgLen,
		Callback:        callback,
	}
	return tcpServer
}

func NewWsserver(
	addr string,
	maxConnNum int,
	pendingWriteNum int,
	maxMsgLen uint32,
	callback Callback,
) *websocketnet.WSServer {
	wsServer := &websocketnet.WSServer{
		Addr:            addr,
		MaxConnNum:      maxConnNum,
		PendingWriteNum: pendingWriteNum,
		MaxMsgLen:       maxMsgLen,
		Callback:        callback,
	}
	return wsServer
}

func NewHttpserver(
	addr string,
	pendingWriteNum int,
	maxMsgLen uint32,
	callback Callback,
) *httpnet.HttpServer {
	httpServer := &httpnet.HttpServer{
		Addr:            addr,
		PendingWriteNum: pendingWriteNum,
		MaxMsgLen:       maxMsgLen,
		Callback:        callback,
	}
	return httpServer
}

func NewTcpclient(addr string) conn.CommConn {
	conn := &tcpnet.TCPClient{
		Addr: addr,
	}

	return conn.Connect()
}

func NewWsclient(addr string) conn.CommConn {
	conn := &websocketnet.WSClient{
		Addr: addr,
	}

	return conn.Connect()
}
