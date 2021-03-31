package main

import (
	"relay/relay"
	"time"
)

var device *relay.MockDevice

func init() {
	device = &relay.MockDevice{
		TCPServerAddress: "0.0.0.0:5000",
		IDM:              "2B",
		IDN:              "08",
		KeepAlive:        1 * time.Second,
	}
}

func main() {
	device.InitTCPClient()
	device.AutoPostDeviceInfo()
	for {
	}
}
