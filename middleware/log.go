package middleware

import (
	"strconv"

	"github.com/fatih/color"
)

var noP = color.New(color.BgHiBlack, color.FgWhite) // 编号输出
var onP = color.New(color.BgBlue, color.FgYellow)   // 闭合输出
var offP = color.New(color.BgRed, color.FgCyan)     // 断开输出
var statePrint = map[string][]func(no int){
	"0": {func(no int) {
		noP.Printf("%d:", no)
		offP.Print("断开")
	}},
	"1": {func(no int) {
		noP.Printf("%d:", no)
		onP.Print("闭合")
	}},
}

// Log 日志输出 用于状态上报
func Log(data interface{}) interface{} {
	switch data.(type) {
	case map[int]interface{}:
		for i, s := range data.(map[int]uint16) {
			for _, p := range statePrint[strconv.Itoa(int(s))] {
				p(i + 1)
			}
		}
	}
	return data
}
