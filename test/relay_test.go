package relay_test

import (
	"fmt"
	"iot-sdk-go/sdk/device"
	"relay/middleware"
	"relay/relay"
	"testing"
	"time"
)

var myRelay *relay.Relay

func init() {
	TCPAddress := "192.168.2.136:9001"     // tcp 监听服务地址
	IOTHubAddress := "39.98.250.155:18100" // iot 平台地址
	productKey := "abdf6b26a399494869c5db5476d1d617fdb5f7d6579fd093ccf78c77ea61e70f"
	deviceName := "relay"
	version := "1.0.0"
	var subDeviceID uint16 = 1
	myRelay = relay.New(TCPAddress, IOTHubAddress, productKey, deviceName, version, subDeviceID, 2*time.Second)
}

func TestInit(t *testing.T) {
	myRelay.Init()
}

func TestAutoPostProperty(t *testing.T) {
	myRelay.Init()
	fns := map[relay.StateType]func() interface{}{
		relay.STATE:      myRelay.GetState,
		relay.THSTATE:    myRelay.GetTHState,
		relay.INPUTSTATE: myRelay.GetInputState,
	}
	myRelay.Use(middleware.Log)
	myRelay.AutoPostProperty(fns)
}

func TestRegisterCommand(t *testing.T) {
	myRelay.Init()
	fns := []device.Command{
		{
			ID: 1,
			Callback: func(params map[int]interface{}) {
				fmt.Println("params:", params)
				// var no = params[0].(uint8)
				// var state = params[1].(uint8)
				// myRelay.SetState(no, state)
			},
		},
	}
	myRelay.RegisterCommand(fns...)
	for {
	}
}

func TestSetState(t *testing.T) {
	myRelay.Init()
	s := []uint8{1, 2, 4, 6, 8}
	fmt.Printf("正在开启第 %v 路\n", s)
	myRelay.SetState(relay.ON, s...)
	time.Sleep(2 * time.Second)
	fmt.Printf("正在关闭第 %v 路\n", s)
	myRelay.SetState(relay.OFF, s...)
	time.Sleep(1 * time.Second)
}
