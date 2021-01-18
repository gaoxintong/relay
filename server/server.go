package server

import (
	"time"
	"iot-sdk-go/sdk/device"
	"iot-sdk-go/sdk/topics"
	"net"
	"relay/server/conn"
	"relay/server/property"
)

var relay *device.Device               // 继电器设备
var tcpAddr = "192.168.2.136:9001"     // tcp 监听服务地址
var httpAddr = "0.0.0.0:8980"          // http 监听服务地址
var iotHubHost = "39.98.250.155:18100" // iot 平台地址
var productKey = "abdf6b26a399494869c5db5476d1d617fdb5f7d6579fd093ccf78c77ea61e70f"
var deviceName = "relay"
var version = "1.0.0"

time.Duration

// Run 启动服务
func Run() {
	// initRelay()
	go startTCPServer()
	for {
	}
	// go startHTTPServer()
}

func initRelay() {
	myTopics := topics.Override(topics.Topics{
		Register: "http://" + iotHubHost + "/v1/devices/registration",
		Login:    "http://" + iotHubHost + "/v1/devices/authentication",
	})
	relay = device.NewDevice(productKey, deviceName, version, device.Topics(myTopics))
	if err := relay.AutoInit(); err != nil {
		panic(err)
	}
}

func startTCPServer() {
	listener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	conn.Conn, err = listener.Accept()
	property.Receive()
}

// func startHTTPServer() {
// 	http.HandleFunc("/cmd", func(rw http.ResponseWriter, req *http.Request) {
// 		q := req.URL.Query()
// 		nos, okno := q["no"]
// 		cmds, okcmd := q["cmd"]
// 		if !okno || !okcmd || len(nos) != 1 || len(cmds) != 1 {
// 			rw.Write([]byte("params parse failed"))
// 			return
// 		}
// 		no := nos[0]
// 		cmdStr := cmds[0]
// 		cmd := strings.ToUpper(cmdStr)
// 		if err := sendCommand(no, cmd); err != nil {
// 			rw.Write([]byte(err.Error()))
// 			return
// 		}
// 		rw.Write([]byte("{\"code\":\"200\"}"))
// 		return
// 	})
// 	http.ListenAndServe(httpAddr, nil)
// }
