package main

import (
	"GinChat/models"
	"GinChat/router"
	"GinChat/utils"
	"github.com/spf13/viper"
	"time"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	InitTimer()
	utils.DB.AutoMigrate(models.UserBasic{})
	utils.DB.AutoMigrate(models.Message{})
	utils.DB.AutoMigrate(models.Contact{})
	utils.DB.AutoMigrate(models.Community{})
	r := router.Router()
	r.Run(viper.GetString("port.server"))
}

// 初始化定时器
func InitTimer() {
	utils.Timer(time.Duration(viper.GetInt("timeout.DelayHeartbeat"))*time.Second, time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second, models.CleanConnection, "")
}
