package main

import (
	"relay/app/relay"
	"time"
)

var device *relay.MockDevice

func init() {
	device = &relay.MockDevice{
		TCPServerAddress: "0.0.0.0:5000",
		IDM:              "10",
		IDN:              "02",
		KeepAlive:        1 * time.Second,
	}
}

func main() {
	device.InitTCPClient()
	device.AutoPostDeviceInfo()
	for {
	}
}
