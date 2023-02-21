package lightning

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"testing"
	"time"

	conn "github.com/yuanxuezhe/lightning/Conn"
)

func Handler(conn conn.CommConn) {
	for {
		//buff, err := network.ReadMsg(conn)
		buff, err := conn.ReadMsg()
		if err != nil {
			break
		}

		conn.WriteMsg([]byte("Hello,Recv msg:" + string(buff)))
		//network.SendMsg(conn, []byte("Hello,Recv msg:"+string(buff)))

		//time.Sleep(1 * time.Millisecond)
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

func Client() {
	conn := NewTcpclient("localhost:8080")
	for i := 0; ; i++ {
		//conn := ynet.NewWsclient("ws://192.168.120.37:8090")
		err := conn.WriteMsg([]byte("YUANSHUAI<==>WANYUAN TCP " + strconv.Itoa(i)))
		if err != nil {
			fmt.Printf("%s", err)
		}
		buff, err := conn.ReadMsg()
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Println(conn.LocalAddr(), "==>", conn.RemoteAddr(), "    ", string(buff))

		// //conn := ynet.NewTcpclient(":8080")
		// conn = ynet.NewWsclient("ws://192.168.2.3:8081")
		// err = conn.WriteMsg([]byte("YUANSHUAI<==>WANYUAN W S " + strconv.Itoa(i)))
		// if err != nil {
		// 	fmt.Printf("%s", err)
		// }
		// buff, err = conn.ReadMsg()
		// if err != nil {
		// 	fmt.Printf("%s", err)
		// }
		// fmt.Println(conn.LocalAddr(), "==>", conn.RemoteAddr(), "    ", string(buff))

		//time.Sleep(1 * time.Second)
	}
}

func TestMain(m *testing.M) {
	// 杀8080端口占用进程
	fmt.Println("kill port:8080")
	processInfo := exec.Command("/bin/bash", "-c", `lsof -i:8080 | awk '{print $2}' | awk  'NR==2{print}'`)
	if pid, err := processInfo.Output(); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(pid))
		processExit := exec.Command("/bin/bash", "-c", `kill -9 `+string(pid))
		if _, err := processExit.Output(); err != nil {
			fmt.Println(err.Error())
		}
	}
	time.Sleep(1 * time.Second)
	fmt.Println("kill port:8080 ok")
	fmt.Println("prepare tcp connection")
	tcpServer := NewTcpserver(
		":8080",
		10,
		1000,
		1000,
		Handler,
	)
	fmt.Println("starting tcp connection")
	tcpServer.Start()
	fmt.Println("start tcp successful")

	fmt.Println("Tcp client test")
	go Client()
	go Client()
	go Client()
	go Client()
	go Client()
	/*
		wsServer := NewWsserver(
			":8081",
			15,
			1000,
			1000,
			Handler,
		)
		wsServer.Start()

		httpServer := NewHttpserver(
			":8082",
			1000,
			1000,
			Handler1,
		)
		httpServer.Start()
	*/
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Printf("System closing down (signal: %v)", sig)

	tcpServer.Close()
}
