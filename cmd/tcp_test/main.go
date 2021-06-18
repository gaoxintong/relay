package main

import (
	"flag"
	"fmt"
	"net"
	"relay/app/gateway"
	"time"
)

func main() {
	TCPAddress := flag.String("tcp-addr", "0.0.0.0:5000", "set tcp address")
	listener, err := net.Listen("tcp4", *TCPAddress)
	if err != nil {
		panic(fmt.Sprintf("tcp listener failed, error: [%v] \n", err))
	}
	fmt.Printf("tcp listener success \n")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[%v] 连接失败，错误：[%v] \n", getNow(), err)
			continue
		}
		data, err := gateway.ReadData(conn, 13)
		if err != nil {
			fmt.Printf("[%v] 读数据失败，错误：[%v] \n", getNow(), err)
			continue
		}
		ID, err := gateway.GetDeviceID(data)
		if err != nil {
			fmt.Printf("[%v] 获取设备 ID 失败，错误：[%v] \n", getNow(), err)
			continue
		}
		fmt.Printf("[%v] 设备：[%v]上线 \n", getNow(), ID)
	}
}

func getNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
