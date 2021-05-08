package main

import (
	"fmt"
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
			KeepAlive:        1 * time.Second,
		})
	}
	return ret
}

func init() {
	devices = makeDevice(6000)
}

func main() {
	fmt.Print(len(devices))
	for _, device := range devices {
		//go func(device *relay.MockDevice) {
		device.InitTCPClient()
		device.AutoPostDeviceInfo()
		//}(device)
	}
	select {}
	//device8808.InitTCPClient()
	//device8808.AutoPostDeviceInfo()
	//for {
	//}
}
