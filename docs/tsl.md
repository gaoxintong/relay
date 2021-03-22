objects：array 属性列表
  no：int 属性编号，对应 PropertyID
  label：string 属性中文名 主要体现在消息转发，在网关层面没有意义
  part：int 冗余设计，暂时没用
  status：array 属性，其中每个对象代表一个数据，对应 Value
    name：string 属性点名称 主要体现在消息转发，在网关层面没有意义，
    value_type：int tlv 的类型

```go
var route uint8 = 1
var value uint8 = 1

{
  PropertyID: 1,
  Value: [route, value]
}
```

commands：array 命令列表
  no：int 命令编号
  name：string 命令中文名
  part：int 冗余设计，暂时没用
  priority：int 冗余设计，暂时没用
  params：array 命令参数列表
    name：string 属性点名称 主要体现在消息转发，在网关层面没有意义
    value_type：int tlv 的类型