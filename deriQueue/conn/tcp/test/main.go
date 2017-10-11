package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8090")
	if err != nil {
		fmt.Println("连接错误")
		fmt.Println(err)
		return
	}
	//a := 1
	defer conn.Close()
	go rev(conn)
	for {
		time.Sleep(1000000000)

		conn.Write([]byte{1, 14, 32})
		// a++
		// if a%1000 == 0 {
		// 	fmt.Println(a)
		// }
	}
}

func rev(conn net.Conn) {
	buff := make([]byte, 50)
	a := 1
	for {
		n, err := conn.Read(buff)

		if err != nil {
			fmt.Println(err)
			continue
		}
		a++
		if a%1000 == 0 {
			fmt.Println(a)
		}
		fmt.Println("接收到数据：" + string(buff[0:n]))
		//fmt.Println(buff[0:n])
	}
}
