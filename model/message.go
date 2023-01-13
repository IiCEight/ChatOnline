package model

import (
	"ChatOnline/wa"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Node struct {
	targetid  uint
	Conn      *websocket.Conn
	DataQueue chan []byte
}

var clientMap map[uint]*Node = make(map[uint]*Node)

var rwLocker sync.RWMutex

func ReadytoChat(userid uint, targetid uint) {
	node := &Node{targetid, nil, make(chan []byte)}
	rwLocker.Lock()
	clientMap[userid] = node
	rwLocker.Unlock()
}

func Chat(userid uint, c *gin.Context) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)
	wa.Checkerr(err)

	rwLocker.Lock()
	clientMap[userid].Conn = conn
	rwLocker.Unlock()
	//读出来
	rwLocker.RLock()
	node := clientMap[userid]
	rwLocker.RUnlock()
	//服务端发送消息
	go sendProc(node)
	//服务端接受消息
	go recvPorc(node.targetid, node)
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			wa.Checkerr(err)
		}
	}
}

func recvPorc(targetid uint, node *Node) {
	for {
		_, msg, err := node.Conn.ReadMessage()
		wa.Checkerr(err)
		fmt.Println("read msg = ", string(msg))
		rwLocker.RLock()
		node, ok := clientMap[targetid]
		rwLocker.RUnlock()
		if ok {
			node.DataQueue <- msg
		}
	}
}
