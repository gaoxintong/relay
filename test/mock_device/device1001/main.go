package main

import (
	"relay/relay"
	"time"
)

var device1001 *relay.MockDevice

func init() {
	device1001 = &relay.MockDevice{
		TCPServerAddress: "192.168.2.136:9001",
		IDM:              "10",
		IDN:              "01",
		KeepAlive:        1 * time.Second,
	}
}

func main() {
	device1001.InitTCPClient()
	device1001.AutoPostDeviceInfo()
	for {
	}
}