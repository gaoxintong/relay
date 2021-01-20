package gateway

import (
	"encoding/binary"
	"fmt"
	"iot-sdk-go/sdk/device"
	"iot-sdk-go/sdk/topics"
	"net"
	"relay/pkg/utils"
	"relay/relay"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Gateway 网关
type Gateway struct {
	Instance *device.Device

	TCPAddress    string
	IOTHubAddress string
	ProductKey    string
	Name          string
	Version       string
	Devices       map[uint16]*relay.Relay
	KeepAlive     time.Duration
}

// New 创建网关
func New(TCPAddress string, IOTHubAddress string, productKey string, gatewayName string, version string) *Gateway {
	return &Gateway{
		TCPAddress:    TCPAddress,
		IOTHubAddress: IOTHubAddress,
		ProductKey:    productKey,
		Name:          gatewayName,
		Version:       version,
		Devices:       make(map[uint16]*relay.Relay),
		KeepAlive:     5 * time.Second,
	}
}

// Run 启动网关服务
func (g *Gateway) Run() error {
	if err := g.initInstance(); err != nil {
		return errors.Wrap(err, "gateway run failed")
	}
	if err := g.initCommand(); err != nil {
		return errors.Wrap(err, "gateway run failed")
	}
	if err := g.startTCPServer(); err != nil {
		return errors.Wrap(err, "gateway run failed")
	}
	return nil
}

// 创建 device 实例
func (g *Gateway) initInstance() error {
	myTopics := topics.Override(topics.Topics{
		Register: "http://" + g.IOTHubAddress + "/v1/devices/registration",
		Login:    "http://" + g.IOTHubAddress + "/v1/devices/authentication",
	})
	deviceInstance := device.NewDevice(g.ProductKey, g.Name, g.Version, device.Topics(myTopics))
	if err := deviceInstance.AutoInit(); err != nil {
		return errors.Wrap(err, "init relay instance failed")
	}
	g.Instance = deviceInstance
	return nil
}

// DeviceOnline 设备上线
func (g *Gateway) DeviceOnline(conn net.Conn, deviceID uint16) {
	device, ok := g.Devices[deviceID]
	if !ok {
		device = relay.New(g.TCPAddress, g.IOTHubAddress, g.ProductKey, g.Name+strconv.Itoa(int(deviceID)), g.Version, deviceID, g.KeepAlive)
		g.Devices[deviceID] = device
	}
	device.Online(g.Instance.PostProperty, conn, relay.DefaultStateTypes)
}

// 注册命令
func (g *Gateway) initCommand() error {
	fns := []device.Command{
		{
			ID:       1,
			Callback: g.dispatchCommand,
		},
	}
	return g.Instance.OnCommand(fns...)
}

// 派遣命令
func (g *Gateway) dispatchCommand(params map[int]interface{}) {
	// 查找子设备
	deviceID, err := makeDeviceID(params[2])
	if err != nil {
		// log
		return
	}
	// 调用子设备的设备状态方法
	if device, ok := g.Devices[deviceID]; ok {
		var no = params[0].([]uint8)[0]
		var state = params[1].([]uint8)[0]
		var stateType relay.StateCMDType
		if state == 1 {
			stateType = relay.ON
		} else if state == 2 {
			stateType = relay.OFF
		} else if state == 3 {
			stateType = relay.DelayedOFF
		}
		device.SetState(stateType, no)
	}
}

// 创建设备 ID
func makeDeviceID(b interface{}) (uint16, error) {
	deviceIDByte, ok := b.([]byte)
	if !ok {
		return 0, errors.New("make device id failed")
	}
	return binary.BigEndian.Uint16(deviceIDByte), nil
}

// 开启 tcp 服务监听
func (g *Gateway) startTCPServer() error {
	listener, err := net.Listen("tcp4", g.TCPAddress)
	if err != nil {
		return errors.Wrap(err, "init tcp server failed")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("连接失败", err)
			// log
			continue
		}
		data, err := readData(conn, 13)
		if err != nil {
			// log
			continue
		}
		ID, err := getDeviceID(data)
		if err != nil {
			// log
			continue
		}
		go g.DeviceOnline(conn, ID)
	}
}

// 读数据
func readData(conn net.Conn, byteOrderLen int) ([]byte, error) {
	data := make([]byte, byteOrderLen)
	if _, err := conn.Read(data); err != nil {
		return data, errors.New("read data failed")
	}
	return data, nil
}

// 获取设备 ID
func getDeviceID(data []byte) (uint16, error) {
	IDM, err := utils.ByteToHex(data[1])
	if err != nil {
		return 0, errors.Wrap(err, "get ID failed")
	}
	IDN, err := utils.ByteToHex(data[2])
	if err != nil {
		return 0, errors.Wrap(err, "get ID failed")
	}
	ID := IDM + IDN
	IDUINT, err := strconv.ParseUint(ID, 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "get ID failed")
	}
	return uint16(IDUINT), nil
}
