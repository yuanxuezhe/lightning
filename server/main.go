package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/solywsh/chatgpt"
	. "github.com/yuanxuezhe/lightning"
	conn "github.com/yuanxuezhe/lightning/Conn"
)

func Handler1(conn conn.CommConn) {
	// The timeout is used to control the situation that the session is in a long and multi session situation.
	// If it is set to 0, there will be no timeout. Note that a single request still has a timeout setting of 30s.
	//chat := chatgpt.New("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", "user_id(not required)", 10*time.Second)
	chat := chatgpt.New("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", "user_id(not required)", 0)
	defer chat.Close()

	for {
		buff, err := conn.ReadMsg()
		if err != nil {
			break
		}

		answer, err := chat.ChatWithContext(string(buff))
		if err != nil {
			fmt.Println(err)
		}

		conn.WriteMsg([]byte(answer))
	}

	//Q: 你认为2022年世界杯的冠军是谁？
	//A: 这个问题很难回答，因为2022年世界杯还没有开始，所以没有人知道冠军是谁。
}

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

func main() {
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

	// 杀8081端口占用进程
	fmt.Println("kill port:8080")
	processInfo = exec.Command("/bin/bash", "-c", `lsof -i:8081 | awk '{print $2}' | awk  'NR==2{print}'`)
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
	fmt.Println("kill port:8081 ok")
	fmt.Println("prepare  connection")

	tcpServer := NewTcpserver(
		":8080",
		10,
		1000,
		1000,
		Handler1,
	)
	tcpServer.Start()

	wsServer := NewWsserver(
		":8081",
		15,
		1000,
		1000,
		Handler,
	)
	wsServer.Start()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Printf("System closing down (signal: %v)", sig)

	tcpServer.Close()
}
