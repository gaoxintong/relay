package main

import (
	"relay/app/relay"
	"time"
)

var device8808 *relay.MockDevice

func init() {
	device8808 = &relay.MockDevice{
		TCPServerAddress: "0.0.0.0:5000",
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
