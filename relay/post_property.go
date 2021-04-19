package relay

import (
	"iot-sdk-go/sdk/device"
	"time"
)

// PropertyIDs 继电器属性 ID 列表
type PropertyIDs struct {
	OutputStateID uint16
	InputStateID  uint16
	TemperatureID uint16
	HumidityID    uint16
}

var relayPropertyIDs = PropertyIDs{
	OutputStateID: 1,
	InputStateID:  2,
	TemperatureID: 3,
	HumidityID:    4,
}

// AutoPostProperty 自动发送属性，会阻塞主线程
func (r *Relay) AutoPostProperty(stateTypes []PropertyType) {
	fns := GetPropertyFnMap{
		OUTPUTSTATE: r.GetOutputState,
		INPUTSTATE:  r.GetInputState,
		TH:          r.GetTH,
	}
	go r.postPropertyLoop(fns, stateTypes)
}

// 循环发送属性
func (r *Relay) postPropertyLoop(fns GetPropertyFnMap, propertyTypes []PropertyType) {
	for {
		select {
		case <-r.closed:
			return
		default:
			time.Sleep(r.keepAlive)
			for _, name := range propertyTypes {
				if _, ok := fns[name]; !ok {
					continue
				}
				property := fns[name]()
				switch name {
				case OUTPUTSTATE:
					r.postOutputState(property)
				case INPUTSTATE:
					r.postInputState(property)
				case TH:
					r.postTH(property)
				}
			}
		}
	}
}

// 发送输出状态
func (r *Relay) postOutputState(property Property) {
	if outputStates, ok := property.(OutputStates); ok {
		// TODO log
		for _, outputState := range outputStates {
			r.PostProperty(relayPropertyIDs.OutputStateID, []interface{}{outputState.Route, outputState.Value})
		}
	}
}

// 发送输入状态
func (r *Relay) postInputState(property interface{}) {
	if inputStates, ok := property.(InputStates); ok {
		// TODO log
		for _, inputState := range inputStates {
			r.PostProperty(relayPropertyIDs.InputStateID, []interface{}{inputState.Route, inputState.Value})
		}
	}
}

// 发送温湿度状态
func (r *Relay) postTH(property interface{}) {
	th, ok := property.(TemperatureAndHumidity)
	if ok {
		r.PostProperty(relayPropertyIDs.TemperatureID, []interface{}{th.Temperature})
		r.PostProperty(relayPropertyIDs.HumidityID, []interface{}{th.Humidity})
	}
}

// PostProperty 发送属性
func (r *Relay) PostProperty(id uint16, value []interface{}) {
	p := device.Property{
		SubDeviceID: r.SubDeviceID,
		PropertyID:  id,
		Value:       value,
	}
	r.Instance.PostProperty(p)
}
