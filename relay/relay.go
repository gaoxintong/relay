package relay

import (
	"encoding/hex"
	"fmt"
	"iot-sdk-go/sdk/device"
	"iot-sdk-go/sdk/topics"
	"net"
	"relay/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const State, THState, InputState = "state", "THState", "InputState"

// Relay 继电器设备
type Relay struct {
	Instance      *device.Device
	TCPAddress    string
	IOTHubAddress string
	ProductKey    string
	DeviceName    string
	Version       string
	SubDeviceID   uint16
	Conn          net.Conn

	middlewares []func(interface{}) interface{}
	state       map[int]uint8
	tHState     map[int]float64
	inputState  map[int]uint8
	d           time.Duration // 发送属性间隔
}

// New 创建继电器实例
func New(TCPAddress string, IOTHubAddress string, productKey string, deviceName string, version string, subDeviceID uint16, duration time.Duration) *Relay {
	return &Relay{
		TCPAddress:    TCPAddress,
		IOTHubAddress: IOTHubAddress,
		ProductKey:    productKey,
		DeviceName:    deviceName,
		Version:       version,
		SubDeviceID:   subDeviceID,
		state:         make(map[int]uint8, 8),
		tHState:       make(map[int]float64, 2),
		inputState:    make(map[int]uint8, 8),
		d:             duration,
	}
}

// Init 初始化资源
func (r *Relay) Init() error {
	if err := r.initInstance(); err != nil {
		return errors.Wrap(err, "init relay failed")
	}
	if err := r.initTCPServer(); err != nil {
		return errors.Wrap(err, "init relay failed")
	}
	if err := r.ReadLoop(13); err != nil {
		return errors.Wrap(err, "init relay failed")
	}
	wfs := []WriteFn{
		{
			fn: r.SendSearchTHStateCMD,
			d:  r.d,
		},
		{
			fn: r.SendSearchInputStateCMD,
			d:  r.d,
		},
	}
	if err := r.WriteLoop(wfs); err != nil {
		return err
	}
	return nil
}

// 创建 device 实例
func (r *Relay) initInstance() error {
	myTopics := topics.Override(topics.Topics{
		Register: "http://" + r.IOTHubAddress + "/v1/devices/registration",
		Login:    "http://" + r.IOTHubAddress + "/v1/devices/authentication",
	})
	relayInstance := device.NewDevice(r.ProductKey, r.DeviceName, r.Version, device.Topics(myTopics))
	if err := relayInstance.AutoInit(); err != nil {
		return errors.Wrap(err, "init relay instance failed")
	}
	r.Instance = relayInstance
	return nil
}

// 开启 tcp 监听
func (r *Relay) initTCPServer() error {
	listener, err := net.Listen("tcp", r.TCPAddress)
	if r.Conn, err = listener.Accept(); err != nil {
		return errors.Wrap(err, "init tcp server failed")
	}
	return nil
}

// RegisterCommand 注册命令
func (r *Relay) RegisterCommand(cmds ...device.Command) error {
	return r.Instance.OnCommand(cmds...)
}

// AutoPostProperty 自动发送属性
func (r *Relay) AutoPostProperty(fns map[string]func() interface{}) {
	go func() {
		for {
			time.Sleep(r.d)
			for name, fn := range fns {
				switch name {
				case State, InputState:
					propertyTyped, ok := fn().(map[int]uint8)
					if ok {
						for id, value := range propertyTyped {
							r.PostProperty(uint16(id), value)
						}
					}
				case THState:
					propertyTyped, ok := fn().(map[int]float64)
					if ok {
						for id, value := range propertyTyped {
							r.PostProperty(uint16(id), value)
						}
					}
				}
			}
		}
	}()
}

// PostProperty 发送属性
func (r *Relay) PostProperty(id uint16, value interface{}) {
	p := device.Property{
		SubDeviceID: r.SubDeviceID,
		PropertyID:  id,
		Value:       []interface{}{value},
	}
	r.Instance.PostProperty(p)
}

// ReadLoop 循环读取数据
func (r *Relay) ReadLoop(byteOrderLen int) error {
	if r.Conn == nil {
		return errors.Wrap(errors.New("not connected"), "read data failed")
	}
	go func() {
		for {
			data := make([]byte, byteOrderLen)
			if _, err := r.Conn.Read(data); err != nil {
				fmt.Println("读取失败")
				// log
				break
			}
			cmd, err := utils.ByteToHex(data[3])
			cmd = strings.ToLower(cmd)
			if err != nil {
				// log
				break
			}
			switch cmd {
			case "aa":
				r.SaveDeviceInfo(data)
			case "1a":
				r.SaveState(data)
			case "1b":
				r.SaveInputState(data)
			case "2a":
				r.SaveTHState(data)
			}
		}
	}()
	return nil
}

// WriteFn 向 conn 写入的方法集合
type WriteFn struct {
	d  time.Duration
	fn func()
}

// WriteLoop 循环发送命令
func (r *Relay) WriteLoop(wfs []WriteFn) error {
	if r.Conn == nil {
		return errors.Wrap(errors.New("not connected"), "write data failed")
	}
	for _, wf := range wfs {
		go func() {
			for {
				time.Sleep(wf.d)
				wf.fn()
			}
		}()
	}
	return nil
}

// Use 使用中间件
func (r *Relay) Use(fns ...func(interface{}) interface{}) {
	for _, fn := range fns {
		r.middlewares = append(r.middlewares, fn)
	}
}

// GetState 获取 8 路状态
func (r *Relay) GetState() interface{} {
	var ret interface{} = r.state
	for _, mw := range r.middlewares {
		ret = mw(ret)
	}
	return ret
}

// GetTHState 获取温湿度状态
func (r *Relay) GetTHState() interface{} {
	var ret interface{} = r.tHState
	for _, mw := range r.middlewares {
		ret = mw(ret)
	}
	return ret
}

// GetInputState 获取 8 路输入状态
func (r *Relay) GetInputState() interface{} {
	var ret interface{} = r.inputState
	for _, mw := range r.middlewares {
		ret = mw(ret)
	}
	return ret
}

// SendSearchTHStateCMD 发送查询温湿度命令
func (r *Relay) SendSearchTHStateCMD() {
	cmd := "A0 01 08 2A 00 00 00 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}

// SendSearchInputStateCMD 发送查询输入状态命令
func (r *Relay) SendSearchInputStateCMD() {
	cmd := "A0 01 08 1B 00 00 00 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}

func (r *Relay) sendCommand(cmd string) error {
	cmd = strings.ReplaceAll(cmd, " ", "")
	cmdByte, err := hex.DecodeString(cmd)
	if err != nil {
		return errors.Wrap(err, "send command failed")
	}
	_, err = r.Conn.Write(cmdByte)
	return errors.Wrap(err, "send command failed")
}

// SaveDeviceInfo 保存设备信息
func (r *Relay) SaveDeviceInfo(data []byte) {
	allStateStr := fmt.Sprintf("%b", data[5])
	if len(allStateStr) == 8 {
		stateMap := make(map[int]uint8, len(allStateStr))
		allStateArr := strings.Split(allStateStr, "")
		for i, state := range allStateArr {
			stateInt, err := strconv.ParseInt(state, 10, 32)
			if err != nil {
				// log
			}
			stateMap[len(allStateArr)-i] = uint8(stateInt)
		}
		r.state = stateMap
	}
}

// SaveState 保存状态
func (r *Relay) SaveState(data []byte) {
	fmt.Println("继电器状态")
	logByte(data)
}

// SaveInputState 保存输入状态
func (r *Relay) SaveInputState(data []byte) {
	allInputStateStr := fmt.Sprintf("%b", data[5])
	if len(allInputStateStr) == 8 {
		inputMap := make(map[int]uint8, len(allInputStateStr))
		allInputStateArr := strings.Split(allInputStateStr, "")
		for i, state := range allInputStateArr {
			inputStateInt, err := strconv.ParseInt(state, 10, 32)
			if err != nil {
				// log
			}
			inputMap[10+len(allInputStateArr)-i] = uint8(inputStateInt)
		}
		r.inputState = inputMap
		fmt.Println("inputMap:", inputMap)
	}
}

// SaveTHState 保存温湿度状态
func (r *Relay) SaveTHState(data []byte) {
	temperaturePositiveOrNegative := fmt.Sprintf("%d", data[4])
	isPositive := false
	if temperaturePositiveOrNegative == "1" {
		isPositive = true
	}
	temperatureInteger := fmt.Sprintf("%d", data[5])
	temperatureDecimal := fmt.Sprintf("%d", data[6])
	temperature, err := stringConcatenationToFloating(temperatureInteger, temperatureDecimal, isPositive)
	if err != nil {
		// log
	}
	humidityInteger := fmt.Sprintf("%d", data[7])
	humidityDecimal := fmt.Sprintf("%d", data[8])
	humidity, err := stringConcatenationToFloating(humidityInteger, humidityDecimal, true)
	if err != nil {
		// log
	}
	THMap := make(map[int]float64, 2)
	THMap[9] = temperature
	THMap[10] = humidity
	r.tHState = THMap
}

// stringConcatenationToFloating 字符串转 float64
func stringConcatenationToFloating(integer string, decimal string, positive bool) (float64, error) {
	str := integer + "." + decimal
	ret, err := strconv.ParseFloat(str, 10)
	if err != nil {
		return 0, err
	}
	if !positive {
		ret -= 2 * ret
	}
	return ret, nil
}

func logByte(data []byte) {
	for _, b := range data {
		fmt.Printf("%x ", b)
	}
}
