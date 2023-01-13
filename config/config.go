package config

import (
	"ChatOnline/wa"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Config struct {
	name     string
	password string
	database string
}

// 加载配置文件
func InitConfig() *Config {
	viper.SetConfigName("config")   //设置配置名
	viper.SetConfigType("yaml")     //告诉配置文件类型
	viper.AddConfigPath("./config") //配置文件的路径
	err := viper.ReadInConfig()     //读取配置文件
	wa.Checkerr(err)
	return &Config{viper.GetString("mysql.name"), viper.GetString("mysql.password"), viper.GetString("mysql.database")}
}

// 链接数据库
func (c *Config) InitMySQL() *gorm.DB {
	// db, err = gorm.Open("mysql", "admin:439956461@(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local")
	dsn := c.name + ":" + c.password + "@(127.0.0.1:3306)/" + c.database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名用单数
		},
	})
	db.AutoMigrate()
	wa.Checkerr(err)
	return db
}
