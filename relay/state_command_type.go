package relay

// StateCMDType 控制状态命令枚举
type StateCMDType string

// ON 闭合
const ON StateCMDType = "01"

// OFF 断开
const OFF StateCMDType = "02"

// DelayedOFF 延迟 3 秒断开
const DelayedOFF StateCMDType = "03"

// Unknown 未知类型
const Unknown StateCMDType = ""
