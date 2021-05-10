package main

import (
	"relay/gateway"
)

func main() {
	// TCPAddress := "192.168.2.136:9001"     // tcp 监听服务地址
	TCPAddress := "0.0.0.0:5000"           // tcp 监听服务地址
	IOTHubAddress := "39.98.250.155:18100" // iot 平台地址
	productKey := "3e67c6e104b920cbc5da9fa8d669d5af36bf672d84cfbce62230addea9de1f6eed89e9ba89ab267aa72c3608a0bd1141"
	deviceName := "relay-test"
	//deviceName := "relay"
	version := "1.0.0"
	g := gateway.New(TCPAddress, IOTHubAddress, productKey, deviceName, version)
	if err := g.Run(); err != nil {
		panic(err)
	}
}
