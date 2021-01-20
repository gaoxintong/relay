package model

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func c() {
	// MySQL 配置信息
	username := "root"     // 账号
	password := "xxxxxxxx" // 密码
	host := "127.0.0.1"    // 地址
	port := 3306           // 端口
	DBname := "gorm1"      // 数据库名称
	timeout := "10s"       // 连接超时，10秒
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, DBname, timeout)
	// Open 连接
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql.")
	}
}
