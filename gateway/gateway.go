package gateway

import (
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
	}
}

// Run 启动网关服务
func (g *Gateway) Run() {
	g.Init()
	g.startTCPServer()
}

// Init 初始化网关
func (g *Gateway) Init() {
	g.initInstance()
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

// AddDevice 添加设备
func (g *Gateway) AddDevice(deviceIDs ...uint16) {
	for _, deviceID := range deviceIDs {
		if _, ok := g.Devices[deviceID]; ok {
			// log
			fmt.Printf("device %v 已存在\n", deviceID)
		}
		g.Devices[deviceID] = relay.New(g.TCPAddress, g.IOTHubAddress, g.ProductKey, g.Name+strconv.Itoa(int(deviceID)), g.Version, deviceID, 5*time.Second)
	}
}

// DeviceOnline 设备上线
func (g *Gateway) DeviceOnline(conn net.Conn, deviceID uint16) error {
	device, ok := g.Devices[deviceID]
	if !ok {
		// log
		return errors.New("device is not registered")
	}
	device.Online(g.Instance.PostProperty, conn, relay.DefaultStateTypes)
	return nil
}

// RegisterCommand 注册命令
func (g *Gateway) RegisterCommand(cmds ...device.Command) error {
	return g.Instance.OnCommand(cmds...)
}

// 派遣命令
func (g *Gateway) DispatchCommand(map[int]interface{}) {

}

// 开启 tcp 服务监听
func (g *Gateway) startTCPServer() {
	listener, err := net.Listen("tcp4", g.TCPAddress)
	if err != nil {
		panic(errors.Wrap(err, "init tcp server failed"))
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
		ID, err := getID(data)
		if err != nil {
			// log
			continue
		}
		go g.DeviceOnline(conn, ID)
	}
}

func readData(conn net.Conn, byteOrderLen int) ([]byte, error) {
	data := make([]byte, byteOrderLen)
	if _, err := conn.Read(data); err != nil {
		return data, errors.New("read data failed")
	}
	return data, nil
}

func getID(data []byte) (uint16, error) {
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
