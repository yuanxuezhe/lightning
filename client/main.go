package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yuanxuezhe/lightning"
)

func main() {
	conn := lightning.NewTcpclient("192.168.1.3:8080")

	for {
		fmt.Printf("Q:")
		// 从标准输入读取输入
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			panic(err)
		}

		err = conn.WriteMsg([]byte(input))
		if err != nil {
			fmt.Printf("%s", err)
		}
		buff, err := conn.ReadMsg()
		if err != nil {
			fmt.Printf("%s", err)
		}

		fmt.Printf("A: %s\n", string(buff))
	}
}
