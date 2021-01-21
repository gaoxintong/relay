package relay

// StateType 状态类型枚举
type StateType string

// STATE 8路开关状态标识符
const STATE StateType = "STATE"

// THSTATE 温湿度状态标识符
const THSTATE StateType = "THSTATE"

// INPUTSTATE 输入状态标识符
const INPUTSTATE StateType = "INPUTSTATE"

// DefaultStateTypes 默认类型列表
var DefaultStateTypes = []StateType{
	STATE, THSTATE, INPUTSTATE,
}
