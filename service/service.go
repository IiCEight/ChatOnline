package service

import (
	"ChatOnline/router"
	"ChatOnline/util"
)

func StartService() {
	util.Initialize()
	router.InitRouter()
}
