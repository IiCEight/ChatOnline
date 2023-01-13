package util

import (
	"ChatOnline/config"

	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() {
	DB = config.InitConfig().InitMySQL()
}
