package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"net/http"
	"SignalServer/utils"
	//"SignalServer/models"
	"SignalServer/models"
)

//type User struct {
//	UserId      string
//	UserType    string
//	GroupTypeId string
//	GroupId     string
//}
//type Content struct {
//	MessageType string
//	Message string
//	UserId      string
//	UserType    string
//	GroupTypeId string
//	GroupId     string
//	SingleSend 	[]string
//	GroupSend bool
//}

type ClientSend struct {
	MessageType  string
	Message  string
	SingleSend 	[]string
	GroupSend bool
}


type WebSocketController struct {
	beego.Controller
}

func leave(group,uuid string)  {
	delete(models.ConnectMap,uuid)
	utils.RedisCache.SRem(group,uuid)
	fmt.Println("开始",models.ConnectMap)
}
func (this *WebSocketController) Join() {
	user := models.User{}
	user.UserId = this.GetString("userId")
	user.UserType = this.GetString("userType")
	user.GroupTypeId = this.GetString("groupTypeId")
	user.GroupId = this.GetString("groupId")
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	uuid :=user.GroupTypeId+"_"+user.GroupId +"_"+ user.UserId
	models.ConnectMap[uuid] = ws
	group := user.GroupTypeId + "_" +user.GroupId
	defer leave(group,uuid)
	utils.RedisCache.SAdd(group,uuid)
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}

		//json转为struct
		receiveMessage := ClientSend{}
		if err := json.Unmarshal([]byte(p), &receiveMessage); err != nil {
			panic(err)
		}
		for index,value := range receiveMessage.SingleSend{
			receiveMessage.SingleSend[index] = user.GroupTypeId+"_"+user.GroupId + "_"+value
		}

		content := models.Content{}
		content.GroupId = user.GroupId
		content.GroupTypeId = user.GroupTypeId
		content.UserType=user.UserType
		content.UserId = user.UserId
		content.GroupSend = receiveMessage.GroupSend
		content.SingleSend = receiveMessage.SingleSend
		content.Message = receiveMessage.Message
		content.MessageType = receiveMessage.MessageType

		b, err := json.Marshal(content)
		if err !=nil{
			panic(err)
		}
		sendMessage := string(b)
		utils.RQ.SendMessage("task_queue",sendMessage)
		// json转为map
		//var receiveMessage map[string]interface{}
		//json.Unmarshal(p, &receiveMessage)
		//newSingleSend := make([]string,0)
		//for _,value := range receiveMessage["singleSend"].([]interface {}) {
		//	if s,ok := value.(string);ok{
		//		sendUser := user.groupTypeId+"_"+user.groupId + "_"+s
		//		newSingleSend = append(newSingleSend,sendUser)
		//	}
		//}
		//receiveMessage["singleSend"] = newSingleSend
		//receiveMessage["userId"]  = user.userId
		//receiveMessage["userType"]  = user.userType
		//receiveMessage["groupTypeId"] = user.groupTypeId
		//receiveMessage["groupId"] = user.groupId
		//fmt.Println(receiveMessage)
		//ws.WriteMessage(websocket.TextMessage,[]byte(sendMessage))
		//fmt.Println(uuid)
	}


}


