package controllers

import (
	"container/list"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"time"
)

type WebSocketController struct {
	beego.Controller
	ws *websocket.Conn
}

type User struct {
	Uid      int
	Username string
	Conn     *websocket.Conn
}

type Msg struct {
	Uid      int
	Username string
	Type     int
	Msg      string
	Time     time.Time
}

var (
	userList  = list.New()
	joinChan  = make(chan User, 10)
	leaveChan = make(chan User, 10)
	uid       = 0
)

const (
	EVENT_LEAVE = iota
	EVENT_JOIN
	EVENT_MSG
)

func join(this *WebSocketController) User {
	username := this.GetString("username")
	ws, _ := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	user := User{Uid: uid, Username: username, Conn: ws}
	joinChan <- user
	for user := userList.Front(); user != nil; user = user.Next() {
		data := user.Value.(User)
		msg := Msg{Uid: data.Uid, Username: data.Username, Type: EVENT_JOIN, Time: time.Now()}
		jdata, err := json.Marshal(msg)
		if err != nil {
			beego.Error("Fail to marshal event:", err)
		}
		ws.WriteMessage(websocket.TextMessage, jdata)
	}
	userList.PushBack(user)
	return user
}
func leave(this *WebSocketController) {
	username := this.GetString("username")
	for user := userList.Front(); user != nil; user = user.Next() {
		name := user.Value.(User).Username
		if name == username {
			leaveChan <- user.Value.(User)
			userList.Remove(user)
		}
	}

}
func send(msg Msg) {
	data, err := json.Marshal(msg)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}
	for user := userList.Front(); user != nil; user = user.Next() {
		conn := user.Value.(User).Conn
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func init() {
	go func() {
		for {
			select {
			case user := <-joinChan:
				msg := Msg{Uid: user.Uid, Username: user.Username, Type: EVENT_JOIN, Time: time.Now()}
				send(msg)
			case user := <-leaveChan:
				msg := Msg{Uid: user.Uid, Username: user.Username, Type: EVENT_LEAVE, Time: time.Now()}
				send(msg)
			}
		}
	}()
}

func (this *WebSocketController) Join() {
	username := this.GetString("username")
	if len(username) == 0 {
		this.Redirect("/", 302)
		return
	}
	uid++
	this.Data["UserName"] = username
	this.Data["Uid"] = uid
	this.TplNames = "room.html"
}

func (this *WebSocketController) Msg() {
	user := join(this)
	defer leave(this)
	for {
		_, p, err := user.Conn.ReadMessage()
		if err != nil {
			return
		}
		msg := Msg{Uid: user.Uid, Username: user.Username, Type: EVENT_MSG, Msg: string(p), Time: time.Now()}
		send(msg)

	}

}

//首页
func (this *WebSocketController) Get() {
	this.TplNames = "join.html"
}
