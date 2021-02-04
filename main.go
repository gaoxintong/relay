package main

import (
	"relay/gateway"
)

func main() {
	TCPAddress := "192.168.2.136:9001"     // tcp 监听服务地址
	IOTHubAddress := "39.98.250.155:18100" // iot 平台地址
	productKey := "abdf6b26a399494869c5db5476d1d617fdb5f7d6579fd093ccf78c77ea61e70f"
	deviceName := "relay"
	version := "1.0.0"
	g := gateway.New(TCPAddress, IOTHubAddress, productKey, deviceName, version)
	if err := g.Run(); err != nil {
		panic(err)
	}
}
