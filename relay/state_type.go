package relay

// PropertyType 属性类型枚举
type PropertyType string

// OUTPUTSTATE 8路开关状态标识符
const OUTPUTSTATE PropertyType = "OUTPUTSTATE"

// TH 温湿度状态标识符
const TH PropertyType = "TH"

// INPUTSTATE 输入状态标识符
const INPUTSTATE PropertyType = "INPUTSTATE"

// DefaultStateTypes 默认类型列表
var DefaultStateTypes = []PropertyType{
	OUTPUTSTATE, TH, INPUTSTATE,
}
