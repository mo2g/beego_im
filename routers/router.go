package routers

import (
	"im/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//首页
	beego.Router("/", &controllers.WebSocketController{})
	beego.Router("/join", &controllers.WebSocketController{}, "get:Join")
	beego.Router("/msg", &controllers.WebSocketController{}, "get:Msg")
}
