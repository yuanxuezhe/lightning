package network

import (
	conn "gitee.com/yuanxuezhe/ynet/Conn"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type WSServer struct {
	Addr            string
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	HTTPTimeout     time.Duration
	//CertFile        string
	//KeyFile         string
	//NewAgent        func(*WSConn) Agent
	ln       net.Listener
	handler  *WSHandler
	Callback func(conn conn.CommConn)
}

type WSHandler struct {
	maxConnNum      int
	pendingWriteNum int
	maxMsgLen       uint32
	//newAgent        func(*WSConn) Agent
	upgrader   websocket.Upgrader
	conns      WebsocketConnSet
	mutexConns sync.Mutex
	wg         sync.WaitGroup
	Callback   func(conn conn.CommConn)
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	conn.SetReadLimit(int64(handler.maxMsgLen))

	handler.wg.Add(1)
	defer handler.wg.Done()

	handler.mutexConns.Lock()
	if handler.conns == nil {
		handler.mutexConns.Unlock()
		conn.Close()
		return
	}
	if len(handler.conns) >= handler.maxConnNum {
		handler.mutexConns.Unlock()
		conn.Close()
		log.Print("too many connections")
		return
	}
	handler.conns[conn] = struct{}{}
	handler.mutexConns.Unlock()

	wsConn := newWSConn(conn, handler.pendingWriteNum, handler.maxMsgLen)

	go func() {
		//agent.Run()
		handler.Callback(wsConn)

		handler.mutexConns.Lock()
		delete(handler.conns, conn)
		handler.mutexConns.Unlock()

		//handler.wg.Done()
	}()
}

func (server *WSServer) Start() {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal("%v", err)
	}

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 100
		log.Printf("invalid MaxConnNum, reset to %v", server.MaxConnNum)
	}
	if server.PendingWriteNum <= 0 {
		server.PendingWriteNum = 100
		log.Printf("invalid PendingWriteNum, reset to %v", server.PendingWriteNum)
	}
	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 4096
		log.Printf("invalid MaxMsgLen, reset to %v", server.MaxMsgLen)
	}
	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second
		log.Printf("invalid HTTPTimeout, reset to %v", server.HTTPTimeout)
	}
	//if server.NewAgent == nil {
	//	log.Fatal("NewAgent must not be nil")
	//}

	//if server.CertFile != "" || server.KeyFile != "" {
	//	config := &tls.Config{}
	//	config.NextProtos = []string{"http/1.1"}
	//
	//	var err error
	//	config.Certificates = make([]tls.Certificate, 1)
	//	config.Certificates[0], err = tls.LoadX509KeyPair(server.CertFile, server.KeyFile)
	//	if err != nil {
	//		log.Fatal("%v", err)
	//	}
	//
	//	ln = tls.NewListener(ln, config)
	//}

	server.ln = ln
	server.handler = &WSHandler{
		maxConnNum:      server.MaxConnNum,
		pendingWriteNum: server.PendingWriteNum,
		maxMsgLen:       server.MaxMsgLen,
		Callback:        server.Callback,
		conns:           make(WebsocketConnSet),
		upgrader: websocket.Upgrader{
			HandshakeTimeout: server.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}

	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.handler,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)
}

func (server *WSServer) Close() {
	server.ln.Close()

	server.handler.mutexConns.Lock()
	for conn := range server.handler.conns {
		conn.Close()
	}
	server.handler.conns = nil
	server.handler.mutexConns.Unlock()

	server.handler.wg.Wait()
}
