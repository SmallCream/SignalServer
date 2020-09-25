package routers

import (
	"SignalServer/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/join/ws",&controllers.WebSocketController{},"get:Join")
	beego.Router("/join",&controllers.WebSocketController{})
}
