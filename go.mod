module relay

go 1.15

require github.com/fatih/color v1.10.0

require (
	github.com/google/pprof v0.0.0-20210413054141-7c2eacd09c8d
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.16.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.11
	iot-sdk-go v0.0.0
)

replace iot-sdk-go => ./pkg/iot-sdk-go
