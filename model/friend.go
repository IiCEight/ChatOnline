package model

import (
	"ChatOnline/util"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type Friend struct {
	gorm.Model
	Username string
	Uid      uint
}

func FindAllFriends(tablename string) []Friend {
	var allfriends []Friend
	util.DB.Debug().Table(tablename).Find(&allfriends)
	return allfriends
}

func FindOneFriendbyUsername(name string, tablename string) *Friend {
	var friend Friend
	util.DB.Debug().Table(tablename).Where("username = ?", name).Find(&friend)
	fmt.Println("friend = ", friend)
	if friend.Username != "" {
		return &friend
	}
	return nil
}

func AddOneFriend(friendname string, username string, uid string, tablename string) *Userinfo {
	friend := FindOneUserbyUsername(friendname)
	if FindOneFriendbyUsername(friend.Username, tablename) != nil {
		return nil
	}
	friendtable := friend.Username + "_" + strconv.FormatUint(uint64(friend.ID), 10)
	//插入当前用户朋友列表
	util.DB.Debug().Table(tablename).Create(&Friend{Username: friend.Username, Uid: friend.ID})
	//更新朋友的朋友列表
	id, _ := strconv.Atoi(uid)
	util.DB.Debug().Table(friendtable).Create(&Friend{Username: username, Uid: uint(id)})
	return friend
}

func DelOneFriend(friendname string, username string, uid string, tablename string) {
	friend := FindOneUserbyUsername(friendname)
	friendtable := friend.Username + "_" + strconv.FormatUint(uint64(friend.ID), 10)
	fmt.Println("friendtabel = ", friendtable)
	fmt.Println("tablename = ", tablename)
	//删除当前用户朋友列表
	util.DB.Debug().Table(tablename).Where("username = ?", friendname).Delete(&Friend{})
	//更新朋友的朋友列表
	util.DB.Debug().Table(friendtable).Where("username = ?", username).Delete(&Friend{})
}
