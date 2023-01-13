package model

import (
	"ChatOnline/util"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
)

// 用户信息
type Userinfo struct {
	gorm.Model
	Username string `form:"username"` //不能缩进妈耶
	Password string `form:"password"`
	Email    string `form:"email"`
}

// 创建用户实例
func InitUserinfo() {

}

func FindOneUserbyID(id uint) (userinfo *Userinfo) {
	util.DB.Where("id = ?", strconv.Itoa(int(id))).Find(&userinfo)
	return userinfo
}

func FindOneUserbyUsername(username string) (userinfo *Userinfo) {
	util.DB.Debug().Where("username = ?", username).Find(&userinfo)
	return userinfo
}

func CreateOneUser(user *Userinfo) {
	fmt.Println(user.Username)
	util.DB.Debug().Create(user)
	util.DB.Debug().Table(user.Username + "_" + strconv.FormatInt(int64(user.ID), 10)).Migrator().CreateTable(&Friend{})
}
