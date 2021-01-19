package middleware

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"

	"relay/relay"
)

var noP = color.New(color.BgHiBlack, color.FgWhite) // 编号输出
var onP = color.New(color.BgBlue, color.FgYellow)   // 闭合输出
var offP = color.New(color.BgRed, color.FgCyan)     // 断开输出
var statePrint = map[string][]func(no int){
	"0": {func(no int) {
		noP.Printf("%d: ", no)
		offP.Print("断开")
	}},
	"1": {func(no int) {
		noP.Printf("%d: ", no)
		onP.Print("闭合")
	}},
}

// 状态Map 转切片 排序
func stateMapToSlice(m map[int]uint8) []uint8 {
	ret := make([]uint8, len(m))
	for k, v := range m {
		ret[k-1] = v
	}
	return ret
}

// 输入状态Map 转切片 排序
func inputStateMapToSlice(m map[int]uint8) []uint8 {
	ret := make([]uint8, len(m))
	for k, v := range m {
		ret[k-11] = v
	}
	return ret
}

// Log 日志输出 用于状态上报
func Log(data relay.Data) relay.Data {
	go func() {
		switch data.StateType {
		case relay.STATE:
			stateLog(data)
		case relay.THSTATE:
			thStateLog(data)
		case relay.INPUTSTATE:
			inputStateLog(data)
		}
	}()
	return data
}

func stateLog(data relay.Data) {
	fmt.Println("继电器状态")
	stateArr := stateMapToSlice(data.Data.(map[int]uint8))
	for i, s := range stateArr {
		for _, p := range statePrint[strconv.Itoa(int(s))] {
			p(i)
		}
	}
	fmt.Println("")
}

func thStateLog(data relay.Data) {
	fmt.Println("")
	th := data.Data.(map[int]float64)
	color.Yellow("温度：%+.2f", th[9]) // 温度
	color.Red("湿度：%.2f", th[10])    // 湿度
	fmt.Println("")
}

func inputStateLog(data relay.Data) {
	fmt.Println("输入状态")
	inputStateArr := inputStateMapToSlice(data.Data.(map[int]uint8))
	for i, s := range inputStateArr {
		fmt.Printf("%v：%v ", i, s)
	}
	fmt.Println("")
}
