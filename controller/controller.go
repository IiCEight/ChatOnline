package controller

import (
	"ChatOnline/model"
	"ChatOnline/wa"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Wrong(c *gin.Context) {
	c.HTML(http.StatusNotFound, "wrong.html", nil)
}

// index
// @Tags 首页
// @Summary 主页
// @Success 200 {string} json{"message"}
// @Router / [get]
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// 展示登录页面
func LoginPage(c *gin.Context) {
	if _, err := c.Cookie("username"); err == nil {
		c.Redirect(http.StatusMovedPermanently, "/home")
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

// 处理登录请求
func Login(c *gin.Context) {
	var user model.Userinfo
	//绑定结构体
	c.ShouldBind(&user)
	if user.Username == "" || user.Password == "" {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"err": "用户名和密码不能为空！",
		})
		return
	}
	resuser := model.FindOneUserbyUsername(user.Username)
	if resuser.Username != user.Username {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"err": "用户名或密码错误!",
		})
	} else {
		c.SetCookie("username", user.Username, 60*60, "/", "127.0.0.1", false, false)
		c.SetCookie("userid", strconv.FormatUint(uint64(resuser.Model.ID), 10), 60*60, "/", "127.0.0.1", false, false)
		fmt.Println("登录成功！")
		c.Redirect(http.StatusMovedPermanently, "/home")
	}
}

// 处理注册请求
func Signup(c *gin.Context) {
	var user model.Userinfo
	//绑定结构体
	c.ShouldBind(&user)
	fmt.Println("Username is ", user.Username)

	if user.Username == "" || user.Password == "" || user.Email == "" { //用户名和密码为空
		c.HTML(http.StatusOK, "login.html", gin.H{
			"err": "用户名和密码不能为空！",
		})
		return
	}
	resuser := model.FindOneUserbyUsername(user.Username)
	if resuser.Username == user.Username && resuser.Username != "" {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"err": "用户名已被占用！",
		})
	} else {
		model.CreateOneUser(&user)
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
}

// 登录中间件
func AlreadyLogin(c *gin.Context) {
	_, err := c.Cookie("username")
	//没有登录
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	c.Next()
}

// 用户主页
func Home(c *gin.Context) {
	// username, err := c.Cookie("username")
	// wa.Checkerr(err)
	c.HTML(http.StatusOK, "home.html", nil)
}

// 处理退出登录请求
func Signout(c *gin.Context) {
	username, err := c.Cookie("username")
	wa.Checkerr(err)
	c.SetCookie("username", username, -1, "/", "127.0.0.1", false, false)
	c.Redirect(http.StatusMovedPermanently, "/login")
}

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

// 升级成websocket协议
func cmdWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	wa.Checkerr(err)
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		fmt.Printf("recv: %s", message)
		err = ws.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}
}

func SendMsg(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	wa.Checkerr(err)
	defer ws.Close()
	msg := "hahhah"
	err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
	wa.Checkerr(err)
	// fmt.Println("jajjaj")
}

func StartChat(c *gin.Context) {
	uid := c.Query("userid")

	userid, _ := strconv.Atoi(uid)
	fmt.Println("userid = ", userid)
	model.Chat(uint(userid), c)
}

func ChatPage(c *gin.Context) {
	friendid := c.Query("friendid")
	uid, _ := c.Cookie("userid")
	userid, _ := strconv.Atoi(uid)
	id, _ := strconv.Atoi(friendid)
	friend := model.FindOneUserbyID(uint(id))

	c.SetCookie("friendname", friend.Username, 60*60, "/", "127.0.0.1", false, false)
	model.ReadytoChat(uint(userid), uint(id))
	c.HTML(http.StatusOK, "index.html", nil)
}

func Finduser(c *gin.Context) {
	friendname := c.Query("username")
	uid, _ := c.Cookie("userid")
	username, _ := c.Cookie("username")
	tablename := username + "_" + uid
	if friendname != "" {
		user := model.AddOneFriend(friendname, username, uid, tablename)
		c.JSON(http.StatusOK, *user)
	} else {
		c.JSON(http.StatusOK, model.FindAllFriends(tablename))
	}
}

func Delfriend(c *gin.Context) {
	friendname := c.Query("friendname")
	uid, _ := c.Cookie("userid")
	username, _ := c.Cookie("username")
	tablename := username + "_" + uid
	model.DelOneFriend(friendname, username, uid, tablename)
}
