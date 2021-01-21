package relay

import (
	"encoding/hex"
	"fmt"
	"iot-sdk-go/sdk/device"
	"net"
	"relay/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Relay 继电器设备
type Relay struct {
	Instance *device.Device

	SubDeviceID uint16
	Conn        net.Conn

	middlewares []func(Data) Data
	state       map[int]uint8
	tHState     map[int]float64
	inputState  map[int]uint8
	keepAlive   time.Duration
}

// New 创建继电器实例
func New(DeviceInstance *device.Device, conn net.Conn, subDeviceID uint16, keepAlive time.Duration) *Relay {
	return &Relay{
		Instance:    DeviceInstance,
		Conn:        conn,
		SubDeviceID: subDeviceID,
		state:       make(map[int]uint8, 8),
		tHState:     make(map[int]float64, 2),
		inputState:  make(map[int]uint8, 8),
		keepAlive:   keepAlive,
	}
}

// Init 初始化资源
func (r *Relay) Init() error {
	// 1. 流读取循环
	if err := r.ReadLoop(13); err != nil {
		return errors.Wrap(err, "init relay failed")
	}
	// 2. 主动询问状态循环
	wfs := []WriteFn{
		{
			fn: r.InquiryTHState,
			d:  r.keepAlive,
		},
		{
			fn: r.InquiryInputState,
			d:  r.keepAlive,
		},
	}
	if err := r.WriteLoop(wfs); err != nil {
		return err
	}
	return nil
}

// AutoPostProperty 自动发送属性，会阻塞主线程
func (r *Relay) AutoPostProperty(stateTypes []StateType) {
	fns := map[StateType]func() interface{}{
		STATE:      r.GetState,
		THSTATE:    r.GetTHState,
		INPUTSTATE: r.GetInputState,
	}
	go r.postPropertyLoop(fns, stateTypes)
	for {
	}
}

// 循环发送属性
func (r *Relay) postPropertyLoop(fns map[StateType]func() interface{}, stateTypes []StateType) {
	for {
		time.Sleep(r.keepAlive)
		for _, name := range stateTypes {
			property := fns[name]()
			switch name {
			case STATE:
				r.postState(property)
			case INPUTSTATE:
				r.postInputState(property)
			case THSTATE:
				r.postTHState(property)
			}
		}
	}
}

// 发送闭合状态
func (r *Relay) postState(property interface{}) {
	propertyTyped, ok := property.(map[int]uint8)
	if ok {
		for id, value := range propertyTyped {
			r.PostProperty(uint16(id), value)
		}
	}
}

// 发送输入状态
func (r *Relay) postInputState(property interface{}) {
	r.postState(property) // 与闭合状态逻辑相同
}

// 发送温湿度状态
func (r *Relay) postTHState(property interface{}) {
	propertyTyped, ok := property.(map[int]float64)
	if ok {
		for id, value := range propertyTyped {
			r.PostProperty(uint16(id), value)
		}
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

// Use 使用中间件
func (r *Relay) Use(fns ...func(Data) Data) {
	for _, fn := range fns {
		r.middlewares = append(r.middlewares, fn)
	}
}

// Data 发送的数据
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

// InquiryTHState 发送查询温湿度命令
func (r *Relay) InquiryTHState() {
	cmd := "A0 01 08 2A 00 00 00 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}

// InquiryInputState 发送查询输入状态命令
func (r *Relay) InquiryInputState() {
	cmd := "A0 01 08 1B 00 00 00 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}

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

// SaveTHState 保存温湿度状态
func (r *Relay) SaveTHState(data []byte) {
	temperature, err := makeTemperature(data[4], data[5], data[6])
	if err != nil {
		fmt.Println(err)
		// log
	}
	humidity, err := makeHumidity(data[7], data[8])
	if err != nil {
		fmt.Println(err)
		// log
	}
	THMap := make(map[int]float64, 2)
	THMap[9] = temperature
	THMap[10] = humidity
	r.tHState = THMap
}

func makeTemperature(positive byte, integer byte, decimal byte) (float64, error) {
	temperaturePositiveOrNegative := fmt.Sprintf("%d", positive)
	isPositive := false
	if temperaturePositiveOrNegative == "1" {
		isPositive = true
	}
	temperatureInteger := fmt.Sprintf("%d", integer)
	temperatureDecimal := fmt.Sprintf("%d", decimal)
	temperature, err := stringToFloat64(temperatureInteger, temperatureDecimal, isPositive)
	if err != nil {
		return 0, errors.Wrap(err, "make temperature failed")
	}
	return temperature, nil
}

func makeHumidity(integer byte, decimal byte) (float64, error) {
	humidityInteger := fmt.Sprintf("%d", integer)
	humidityDecimal := fmt.Sprintf("%d", decimal)
	humidity, err := stringToFloat64(humidityInteger, humidityDecimal, true)
	if err != nil {
		return 0, errors.Wrap(err, "make humidity failed")
	}
	return humidity, nil
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

// Online 上线
func (r *Relay) Online(stateTypes []StateType) error {
	if err := r.Init(); err != nil {
		return err
	}
	r.AutoPostProperty(stateTypes)
	return nil
}
