package main

import (
	"relay/pkg/convcode"
	"relay/relay"
	"time"
)

var devices []*relay.MockDevice

func makeDevice(num uint16) []*relay.MockDevice {
	var ret []*relay.MockDevice
	var i uint16 = 1000

	for ; i < num+1000; i++ {
		id := convcode.Dec2Hex(int(i))
		ret = append(ret, &relay.MockDevice{
			TCPServerAddress: "0.0.0.0:5000",
			IDM:              id[0:1],
			IDN:              id[2:3],
			KeepAlive:        2 * time.Second,
		})
	}
	return ret
}

func init() {
	devices = makeDevice(20)
}

func main() {
	for _, device := range devices {
		go func(device *relay.MockDevice) {
			device.InitTCPClient()
			device.AutoPostDeviceInfo()
		}(device)
	}
	select {}
	//device8808.InitTCPClient()
	//device8808.AutoPostDeviceInfo()
	//for {
	//}
}
