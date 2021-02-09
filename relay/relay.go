package relay

import (
	"iot-sdk-go/sdk/device"
	"net"
	"time"

	"github.com/pkg/errors"
)

// Relay 继电器设备
type Relay struct {
	Instance *device.Device

	SubDeviceID uint16
	Conn        net.Conn

	middlewares []func(Data) Data
	outputState OutputStates
	inputState  InputStates
	th          TemperatureAndHumidity
	keepAlive   time.Duration
}

// OutputStates 输出状态集合
type OutputStates []OutputState

// OutputState 输出状态
type OutputState struct {
	Route uint8
	Value uint8
}

// InputStates 输入状态集合
type InputStates []InputState

// InputState 输入状态
type InputState struct {
	Route uint8
	Value uint8
}

// TemperatureAndHumidity 温湿度
type TemperatureAndHumidity struct {
	Temperature float64 // 温度
	Humidity    float64 //湿度
}

// GetPropertyFnMap 根据属性类型获取属性的方法集合
type GetPropertyFnMap map[PropertyType]GetPropertyFn

// GetPropertyFn 不同属性类型的获取属性方法
type GetPropertyFn func() Property

// Property 属性
type Property interface{}

// New 创建继电器实例
func New(DeviceInstance *device.Device, conn net.Conn, subDeviceID uint16, keepAlive time.Duration) *Relay {
	return &Relay{
		Instance:    DeviceInstance,
		Conn:        conn,
		SubDeviceID: subDeviceID,
		outputState: OutputStates{},
		inputState:  InputStates{},
		th:          TemperatureAndHumidity{},
		keepAlive:   keepAlive,
	}
}

// Init 初始化资源
func (r *Relay) Init() error {
	// 流读取循环
	if err := r.ReadLoop(13); err != nil {
		return errors.Wrap(err, "init relay failed")
	}
	// 主动询问状态循环
	wfs := []WriteFn{
		{
			fn: r.InquiryTH,
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

// Use 使用中间件
func (r *Relay) Use(fns ...func(Data) Data) {
	for _, fn := range fns {
		r.middlewares = append(r.middlewares, fn)
	}
}

// Online 上线
func (r *Relay) Online(stateTypes []PropertyType) error {
	if err := r.Init(); err != nil {
		return err
	}
	r.AutoPostProperty(stateTypes)
	return nil
}
