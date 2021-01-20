module relay

go 1.15

require github.com/fatih/color v1.10.0

require (
	github.com/pkg/errors v0.9.1
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.11
	iot-sdk-go v0.0.0
)

replace iot-sdk-go => ./pkg/iot-sdk-go
