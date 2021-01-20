package main

import (
	"relay/relay"
	"time"
)

var device8808 *relay.MockDevice

func init() {
	device8808 = &relay.MockDevice{
		TCPServerAddress: "192.168.2.136:9001",
		IDM:              "88",
		IDN:              "08",
		KeepAlive:        1 * time.Second,
	}
}

func main() {
	device8808.InitTCPClient()
	device8808.AutoPostDeviceInfo()
	for {
	}
}
