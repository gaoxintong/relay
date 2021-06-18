package gateway

import (
	"encoding/binary"
	"fmt"
	"iot-sdk-go/sdk/device"
	"iot-sdk-go/sdk/topics"
	"net"
	relay2 "relay/app/relay"
	"relay/pkg/convcode"
	"relay/pkg/lru"
	"relay/pkg/utils"
	"time"

	_ "net/http/pprof"

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
	Devices       Devices
	KeepAlive     time.Duration
}

type Devices map[uint16]*Device

type Device struct {
	*relay2.Relay
	*lru.LRUCache
}

type DataRecord struct {
	relay2.Data
	Time string
}

// New 创建网关
func New(TCPAddress string, IOTHubAddress string, productKey string, gatewayName string, version string) *Gateway {
	return &Gateway{
		TCPAddress:    TCPAddress,
		IOTHubAddress: IOTHubAddress,
		ProductKey:    productKey,
		Name:          gatewayName,
		Version:       version,
		Devices:       make(map[uint16]*Device, 100),
		KeepAlive:     5 * time.Second,
	}
}

// Run 启动网关服务
func (g *Gateway) Run() error {
	if err := g.initInstance(); err != nil {
		return errors.Wrap(err, "gateway: init instance failed")
	}
	if err := g.initCommand(); err != nil {
		return errors.Wrap(err, "gateway: init command failed")
	}
	go g.startHTTPServer()
	if err := g.startTCPServer(); err != nil {
		return errors.Wrap(err, "gateway: start tcp server failed")
	}
	return nil
}

// 创建 device 实例
func (g *Gateway) initInstance() error {
	myTopics := topics.Override(topics.Topics{
		Register: "http://" + g.IOTHubAddress + "/v1/devices/registration",
		Login:    "http://" + g.IOTHubAddress + "/v1/devices/authentication",
	})
	deviceInstance := device.New(g.ProductKey, g.Name, g.Version, device.Topics(myTopics))
	if err := deviceInstance.AutoInit(); err != nil {
		return errors.Wrap(err, "init relay instance failed")
	}
	g.Instance = deviceInstance
	return nil
}

// DeviceOnline 设备上线
func (g *Gateway) DeviceOnline(conn net.Conn, deviceID uint16) {
	device := g.getOrCreateRelay(conn, deviceID)
	device.Online(relay2.DefaultStateTypes)
}

// 获取继电器设备实例 不存在则创建
func (g *Gateway) getOrCreateRelay(conn net.Conn, deviceID uint16) *relay2.Relay {
	device, ok := g.Devices[deviceID]
	if !ok {
		device = &Device{g.makeRelay(conn, deviceID), lru.New(20)}
		g.Devices[deviceID] = device
	}
	return device.Relay
}

// 创建继电器设备实例
func (g *Gateway) makeRelay(conn net.Conn, deviceID uint16) *relay2.Relay {
	offlineCb := func(relay *relay2.Relay) {
		delete(g.Devices, uint16(relay.SubDeviceID))
	}

	return relay2.New(g.Instance,
		conn, deviceID,
		g.KeepAlive,
		relay2.OfflineCallback(offlineCb),
		relay2.Middlewares(PropertyLog(g.Devices)),
	)
}

// 注册命令
func (g *Gateway) initCommand() error {
	fns := []device.Command{
		{
			ID:       1, // 1 表示控制输入开关
			Callback: g.dispatchSetStateCommand,
		},
	}
	return g.Instance.OnCommand(fns...)
}

func (g *Gateway) dispatchSetStateCommand(params map[int]interface{}) {
	// 查找子设备
	deviceID, err := makeDeviceID(params[2])
	if err != nil {
		// log
		return
	}
	// 调用子设备的设备状态方法
	if device, ok := g.Devices[deviceID]; ok {
		no := params[0].([]uint8)[0]
		state := params[1].([]uint8)[0]
		stateType := getDeviceState(state)
		if stateType != relay2.Unknown {
			device.SetState(stateType, no)
		}
	}
}

func getDeviceState(state uint8) relay2.StateCMDType {
	if state == 1 {
		return relay2.ON
	}
	if state == 2 {
		return relay2.OFF
	}
	if state == 3 {
		return relay2.DelayedOFF
	}
	return relay2.Unknown
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
		data, err := ReadData(conn, 13)
		if err != nil {
			// log
			continue
		}
		ID, err := GetDeviceID(data)
		if err != nil {
			// log
			continue
		}
		go g.DeviceOnline(conn, ID)
	}
}

// 读数据
func ReadData(conn net.Conn, byteOrderLen int) ([]byte, error) {
	data := make([]byte, byteOrderLen)
	if _, err := conn.Read(data); err != nil {
		return data, errors.New("read data failed")
	}
	return data, nil
}

// 获取设备 ID
func GetDeviceID(data []byte) (uint16, error) {
	IDM, err := utils.ByteToHex(data[1])
	if err != nil {
		return 0, errors.Wrap(err, "get ID failed")
	}
	IDN, err := utils.ByteToHex(data[2])
	if err != nil {
		return 0, errors.Wrap(err, "get ID failed")
	}
	IDStr := IDM + IDN
	IDUint16, err := convcode.Hex2Dec(IDStr)
	if err != nil {
		return 0, errors.Wrap(err, "get ID failed")
	}
	return uint16(IDUint16), nil
}
