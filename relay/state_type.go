package relay

// PropertyType 属性类型枚举
type PropertyType string

// STATE 8路开关状态标识符
const STATE PropertyType = "STATE"

// THSTATE 温湿度状态标识符
const THSTATE PropertyType = "THSTATE"

// INPUTSTATE 输入状态标识符
const INPUTSTATE PropertyType = "INPUTSTATE"

// DefaultStateTypes 默认类型列表
var DefaultStateTypes = []PropertyType{
	STATE, THSTATE, INPUTSTATE,
}
