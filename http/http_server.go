package network

import (
	"log"
	"net"
	"net/http"
	"time"

	conn "github.com/yuanxuezhe/lightning/Conn"
)

type HttpServer struct {
	Addr            string
	PendingWriteNum int
	MaxMsgLen       uint32
	HTTPTimeout     time.Duration
	//CertFile        string
	//KeyFile         string
	//NewAgent        func(*HttpConn) Agent
	ln       net.Listener
	handler  *HttpHandler
	Callback func(conn conn.CommConn)
}

type HttpHandler struct {
	//maxConnNum      int
	pendingWriteNum int
	maxMsgLen       uint32
	//upgrader   websocket.Upgrader
	//mutexConns sync.Mutex
	//wg         sync.WaitGroup
	Callback func(conn conn.CommConn)
}

func (handler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpConn := newHttpConn(w, r, handler.pendingWriteNum, handler.maxMsgLen)

	go func() {
		handler.Callback(httpConn)
	}()
}

func (server *HttpServer) Start() {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal("%v", err)
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
	server.handler = &HttpHandler{
		Callback: server.Callback,
	}

	httpServer := &http.Server{
		Addr:    server.Addr,
		Handler: server.handler,
	}

	go httpServer.Serve(ln)
}

func (server *HttpServer) Close() {
}
