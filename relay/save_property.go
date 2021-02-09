package relay

import (
	"fmt"
	"relay/pkg/utils"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// SaveOutputState 保存输出状态
func (r *Relay) SaveOutputState(data []byte) {
	stateStr := utils.ByteToBinaryString(data[5]) // 二进制字符串
	stateStrs := strings.Split(stateStr, "")      // 字符串数组
	states, err := makeOutputStates(stateStrs)
	if err != nil {
		// TODO log
	}
	r.outputState = states
}

func makeOutputStates(stateStrs []string) (OutputStates, error) {
	states := OutputStates{}
	for i, state := range stateStrs {
		stateInt, err := strconv.ParseInt(state, 10, 8)
		if err != nil {
			return nil, err
		}
		states = append(states, OutputState{
			Route: uint8(len(stateStrs) - i),
			Value: uint8(stateInt),
		})
	}
	return states, nil
}

// SaveInputState 保存输入状态
func (r *Relay) SaveInputState(data []byte) {
	stateStr := utils.ByteToBinaryString(data[5]) // 二进制字符串
	stateStrs := strings.Split(stateStr, "")      // 字符串数组
	states, err := makeInputStates(stateStrs)
	if err != nil {
		// TODO log
	}
	r.inputState = states
}

func makeInputStates(stateStrs []string) (InputStates, error) {
	states := InputStates{}
	for i, state := range stateStrs {
		stateInt, err := strconv.ParseInt(state, 10, 8)
		if err != nil {
			return nil, err
		}
		states = append(states, InputState{
			Route: uint8(len(stateStrs) - i),
			Value: uint8(stateInt),
		})
	}
	return states, nil
}

// SaveTH 保存温湿度状态
func (r *Relay) SaveTH(data []byte) {
	temperature, err := makeTemperature(data[4], data[5], data[6])
	if err != nil {
		fmt.Println(err)
		// TODO log
	}
	humidity, err := makeHumidity(data[7], data[8])
	if err != nil {
		fmt.Println(err)
		// TODO log
	}
	th := TemperatureAndHumidity{
		Temperature: temperature,
		Humidity:    humidity,
	}
	r.th = th
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
