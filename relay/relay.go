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

// StateType 状态类型枚举
type StateType string

// STATE 8路开关状态标识符
const STATE StateType = "STATE"

// THSTATE 温湿度状态标识符
const THSTATE StateType = "THSTATE"

// INPUTSTATE 输入状态标识符
const INPUTSTATE StateType = "INPUTSTATE"

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

	middlewares []func(Data) Data
	state       map[int]uint8
	tHState     map[int]float64
	inputState  map[int]uint8
	d           time.Duration // 上报属性间隔时间
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
	// 1. 创建 device 实例
	if err := r.initInstance(); err != nil {
		return errors.Wrap(err, "init relay failed")
	}
	// 2. 开启 TCP监听
	if err := r.initTCPServer(); err != nil {
		return errors.Wrap(err, "init relay failed")
	}
	// 3. 流读取循环
	if err := r.ReadLoop(13); err != nil {
		return errors.Wrap(err, "init relay failed")
	}
	// 4. 主动询问状态循环
	wfs := []WriteFn{
		{
			fn: r.SearchTHState,
			d:  r.d,
		},
		{
			fn: r.SearchInputState,
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

// AutoPostProperty 自动发送属性，会阻塞主线程
func (r *Relay) AutoPostProperty(stateTypes []StateType) {
	fns := map[StateType]func() interface{}{
		STATE:      r.GetState,
		THSTATE:    r.GetTHState,
		INPUTSTATE: r.GetInputState,
	}
	go func() {
		for {
			time.Sleep(r.d)
			for _, name := range stateTypes {
				switch name {
				case STATE, INPUTSTATE:
					propertyTyped, ok := fns[name]().(map[int]uint8)
					if ok {
						for id, value := range propertyTyped {
							r.PostProperty(uint16(id), value)
						}
					}
				case THSTATE:
					propertyTyped, ok := fns[name]().(map[int]float64)
					if ok {
						for id, value := range propertyTyped {
							r.PostProperty(uint16(id), value)
						}
					}
				}
			}
		}
	}()
	for {
	}
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

// ReadLoop 开启一个协程，从连接中循环读取数据
func (r *Relay) ReadLoop(byteOrderLen int) error {
	if r.Conn == nil {
		return errors.Wrap(errors.New("not connected"), "read data failed")
	}
	go func() {
		for {
			data := make([]byte, byteOrderLen)
			if _, err := r.Conn.Read(data); err != nil {
				fmt.Println("读取失败", err)
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

// WriteLoop 开启N个协程，向连接循环发送命令
func (r *Relay) WriteLoop(wfs []WriteFn) error {
	if r.Conn == nil {
		return errors.Wrap(errors.New("not connected"), "write data failed")
	}
	for _, wf := range wfs {
		go func(wf WriteFn) {
			for {
				time.Sleep(wf.d)
				wf.fn()
			}
		}(wf)
	}
	return nil
}

// Use 使用中间件
func (r *Relay) Use(fns ...func(Data) Data) {
	for _, fn := range fns {
		r.middlewares = append(r.middlewares, fn)
	}
}

type Data struct {
	StateType StateType
	Data      interface{}
}

// 获取数据
func (r *Relay) getData(data Data) interface{} {
	for _, mw := range r.middlewares {
		data = mw(data)
	}
	return data.Data
}

// GetState 获取 8 路状态
func (r *Relay) GetState() interface{} {
	return r.getData(Data{StateType: STATE, Data: r.state})
}

// GetTHState 获取温湿度状态
func (r *Relay) GetTHState() interface{} {
	return r.getData(Data{StateType: THSTATE, Data: r.tHState})
}

// GetInputState 获取 8 路输入状态
func (r *Relay) GetInputState() interface{} {
	return r.getData(Data{StateType: INPUTSTATE, Data: r.inputState})
}

// SearchTHState 发送查询温湿度命令
func (r *Relay) SearchTHState() {
	cmd := "A0 01 08 2A 00 00 00 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}

// SearchInputState 发送查询输入状态命令
func (r *Relay) SearchInputState() {
	cmd := "A0 01 08 1B 00 00 00 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}

// StateCMDType 控制状态命令枚举
type StateCMDType string

// ON 闭合
const ON StateCMDType = "01"

// OFF 断开
const OFF StateCMDType = "02"

// DelayedOFF 延迟 3 秒断开
const DelayedOFF StateCMDType = "03"

// SetState 设置状态
func (r *Relay) SetState(state StateCMDType, nos ...uint8) {
	no, err := getNoHex(nos...)
	if err != nil {
		// log
	}
	cmd := "A0 01 08 2B 00 " + no + " " + string(state) + " 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}

// 获取选择编号 16 进制字符串
func getNoHex(nos ...uint8) (string, error) {
	stateArr := strings.Split("00000000", "")
	for _, no := range nos {
		stateArr[7-(no-1)] = "1"
	}
	stateStr := strings.Join(stateArr, "")
	return utils.BinaryToHex(stateStr)
}

// 发送命令
func (r *Relay) sendCommand(cmd string) error {
	cmdByte, err := commandFormatter(cmd)
	if err != nil {
		return errors.Wrap(err, "send command failed")
	}
	if _, err = r.Conn.Write(cmdByte); err != nil {
		return errors.Wrap(err, "send command failed")
	}
	return nil
}

// 格式化命令
func commandFormatter(originalCommand string) ([]byte, error) {
	originalCommand = strings.ReplaceAll(originalCommand, " ", "")
	return hex.DecodeString(originalCommand)
}

// SaveState 保存状态
func (r *Relay) SaveState(data []byte) {
	allStateStr := fmt.Sprintf("%b", data[5])
	stateMap := make(map[int]uint8, len(allStateStr))
	if len(allStateStr) == 8 {
		allStateArr := strings.Split(allStateStr, "")
		for i, state := range allStateArr {
			stateInt, err := strconv.ParseInt(state, 10, 32)
			if err != nil {
				// log
			}
			stateMap[len(allStateArr)-i] = uint8(stateInt)
		}
		r.state = stateMap
	} else if allStateStr == "0" {
		for i := 1; i <= 8; i++ {
			stateMap[i] = 0
		}
		r.state = stateMap
	}
}

// SaveInputState 保存输入状态
func (r *Relay) SaveInputState(data []byte) {
	allInputStateStr := utils.ByteToBinary(data[4])
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
	temperature, err := stringToFloat64(temperatureInteger, temperatureDecimal, isPositive)
	if err != nil {
		// log
	}
	humidityInteger := fmt.Sprintf("%d", data[7])
	humidityDecimal := fmt.Sprintf("%d", data[8])
	humidity, err := stringToFloat64(humidityInteger, humidityDecimal, true)
	if err != nil {
		// log
	}
	THMap := make(map[int]float64, 2)
	THMap[9] = temperature
	THMap[10] = humidity
	r.tHState = THMap
}

// stringToFloat64 字符串转 float64
func stringToFloat64(integer string, decimal string, positive bool) (float64, error) {
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
