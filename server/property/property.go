package property

import (
	"encoding/hex"
	"fmt"
	"relay/pkg/utils"
	"relay/server/command"
	"relay/server/conn"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/fatih/color"
)

var byteOrderLen = 13                               //字节序长度
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

// Receive 被动接收状态
func Receive() {
	defer func() {
		if conn.Conn != nil {
			conn.Conn.Close()
		}
	}()
	command.SendCommand("01", "OFF")
	// [160 1 8 43 0 1 2 0 0 0 0 0 167]
	// [160 1 8 43 0 1 0 0 0 0 0 0 167]
	for {
		time.Sleep(5 * time.Second)
		go func() {
			fmt.Println()
			// getState()
			GetTemperature()
			// onP.Printf("温湿度：%v", th)
		}()
	}
}

// 被动获取 8 路状态
func getState() {
	data, err := conn.ReadData(byteOrderLen)
	if err != nil {

	}
	for i, d := range data {
		if i == 5 {
			allStateStr := utils.ByteToBinary(d)
			allState := strings.Split(allStateStr, "")
			for i, s := range allState {
				for _, p := range statePrint[s] {
					p(i + 1)
				}
			}
		}
	}
}

func get(cmd string) ([]byte, error) {
	if err := conn.CheckConnState(); err != nil {
		return nil, errors.Wrap(err, "get state failed")
	}
	cmdByte, err := getCmdByte(cmd)
	fmt.Println("执行命令:", cmd)
	if err != nil {
		return nil, errors.Wrap(err, "exec "+cmd+" failed")
	}
	conn.CmdConn.Write(cmdByte)
	return conn.ReadData(byteOrderLen)
}

func getCmdByte(cmd string) ([]byte, error) {
	cmd = strings.ReplaceAll(cmd, " ", "") // 去除空格
	return hex.DecodeString(cmd)           // 字符串 转 byte
}

// GetState 获取设备信息
func GetDeviceInfo() (string, error) {
	// cmd := "A0 01 08 AA 00 00 00 00 00 00 00 A7"
	// return get(cmd)
	return "", nil
}

// GetInputState 获取输入状态
func GetInputState() (string, error) {
	cmd := "A0 01 08 1B 00 00 00 00 00 00 00 A7"
	data, _ := get(cmd)
	for _, v := range data {
		fmt.Printf("%X ", v)
	}
	// return get(cmd)
	return "", nil
}

// GetTemperature 主动获取温度
func GetTemperature() (string, error) {
	cmd := "A0 01 08 2A 00 00 00 00 00 00 00 00 A7"
	data, _ := get(cmd)
	fmt.Println("温度：")
	for _, v := range data {
		fmt.Printf("%X ", v)
	}
	return "", nil
}

// GetHumidity 主动获取湿度
func GetHumidity() (string, error) {
	// cmd := "A0 01 08 2A 00 00 00 00 00 00 00 A7"
	// return get(cmd)
	return "", nil
}
