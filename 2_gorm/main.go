package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DailyMetrics struct {
	gorm.Model
	ShowTime    time.Time
	MetricValue string `gorm:"not null"`
	MetricName  string `gorm:"not null"`
}

type OuterData struct {
	gorm.Model
	ShowTime time.Time
	Value    string `gorm:"not null"`
	Name     string `gorm:"not null"`
}

func main() {
	dsn := "root:root1234@tcp(sh-cdb-9xiziyvi.sql.tencentcdb.com:63879)/testgorm?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("连接数据库成功")

	// 迁移 schema
	db.AutoMigrate(&DailyMetrics{})
	db.AutoMigrate(&OuterData{})

	// 创建数据
	db.Create(&OuterData{ShowTime: time.Now(), Value: "666", Name: "wcxiao"})
	ds := []*DailyMetrics{
		{
			ShowTime:    time.Now().AddDate(0, 0, -1),
			MetricValue: "gtmd",
			MetricName:  "yesterday",
		},
		{
			ShowTime:    time.Now(),
			MetricValue: "6666",
			MetricName:  "today",
		},
	}
	db.CreateInBatches(ds, 1)
}
