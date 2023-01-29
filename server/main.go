package main

import (
	"fmt"
	"gitee.com/yuanxuezhe/ynet"
	conn "gitee.com/yuanxuezhe/ynet/Conn"
	"os"
	"os/signal"
	"time"
)

//func Handler1(conn net.Conn) {
//	for {
//
//		msgs, err := network.ReadMsgs(conn)
//		if err != nil {
//			break
//		}
//		for _, msg := range msgs {
//			fmt.Println(msg)
//			network.SendMsg(conn,[]byte("Hello,Recv msg:" + msg))
//		}
//
//		time.Sleep(10 * time.Millisecond)
//	}
//}

func Handler(conn conn.CommConn) {
	for {
		//buff, err := network.ReadMsg(conn)
		buff, err := conn.ReadMsg()
		if err != nil {
			break
		}

		conn.WriteMsg([]byte("Hello,Recv msg:" + string(buff)))
		//network.SendMsg(conn, []byte("Hello,Recv msg:"+string(buff)))

		time.Sleep(1 * time.Millisecond)
	}
}

func Handler1(conn conn.CommConn) {
	//buff, err := network.ReadMsg(conn)
	buff, err := conn.ReadMsg()
	fmt.Println("1111111111111")
	if err != nil {
		return
	}
	fmt.Println("22222222222222")
	fmt.Println("Recv:", string(buff))

	conn.WriteMsg([]byte("Hello,Recv msg:" + string(buff)))
}

func main() {
	tcpServer := ynet.NewTcpserver(
		":8080",
		10,
		1000,
		1000,
		Handler,
	)
	tcpServer.Start()

	wsServer := ynet.NewWsserver(
		":8081",
		15,
		1000,
		1000,
		Handler,
	)
	wsServer.Start()

	httpServer := ynet.NewHttpserver(
		":8082",
		1000,
		1000,
		Handler1,
	)
	httpServer.Start()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Printf("System closing down (signal: %v)", sig)

	tcpServer.Close()
}
