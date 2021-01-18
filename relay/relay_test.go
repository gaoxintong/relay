package relay

import (
	"fmt"
	"iot-sdk-go/sdk/device"
	"relay/middleware"
	"testing"
	"time"
)

var myRelay *Relay

func init() {
	TCPAddress := "192.168.2.136:9001"     // tcp 监听服务地址
	IOTHubAddress := "39.98.250.155:18100" // iot 平台地址
	productKey := "abdf6b26a399494869c5db5476d1d617fdb5f7d6579fd093ccf78c77ea61e70f"
	deviceName := "relay"
	version := "1.0.0"
	var subDeviceID uint16 = 1
	myRelay = New(TCPAddress, IOTHubAddress, productKey, deviceName, version, subDeviceID, 2*time.Second)
}

func TestInit(t *testing.T) {
	myRelay.Init()
	for {
	}
}

func TestAutoPostProperty(t *testing.T) {
	myRelay.Init()
	fns := map[string]func() interface{}{
		State:      myRelay.GetState,
		THState:    myRelay.GetTHState,
		InputState: myRelay.GetInputState,
	}
	myRelay.Use(middleware.Log)
	myRelay.AutoPostProperty(fns)
	for {
	}
}

func TestRegisterCommand(t *testing.T) {
	myRelay.Init()
	fns := []device.Command{
		{
			ID: 1,
			Callback: func(params map[int]interface{}) {
				fmt.Println("params:", params)
			},
		},
	}
	myRelay.RegisterCommand(fns...)
	for {
	}
}
