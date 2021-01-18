module relay

go 1.15

require github.com/fatih/color v1.10.0

require (
	github.com/pkg/errors v0.9.1
	iot-sdk-go v0.0.0
)

replace iot-sdk-go => ./pkg/iot-sdk-go
