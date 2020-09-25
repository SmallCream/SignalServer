package utils

import (
	"SignalServer/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

var RQ = Rabbitmq{}

type Rabbitmq struct {
	ConnectUrl string
	Conn *amqp.Connection
	Ch *amqp.Channel
}

func (r *Rabbitmq) Connect()  {
	var err error
	r.Conn,err = amqp.Dial(r.ConnectUrl)
	if err!=nil{
		fmt.Println("Rabbitmq创建连接失败")
	}
	r.Ch, err = r.Conn.Channel()
	if err != nil{
		fmt.Println("Rabbitmq创建管道失败")
	}
}

func (r *Rabbitmq) Receive(queue_name string) {
	mssages,err:=r.Ch.Consume(queue_name,"",false,false,false,false,nil)
	if err != nil{
		fmt.Println("rabbitmq接受消息失败")
	}
	go func() {
		for msg := range mssages{

			content := models.Content{}
			if err := json.Unmarshal(msg.Body, &content); err != nil {
				fmt.Println("rabbitmq消息转struct出错",err)
			}
			clientReceive := models.ClientReceive{}
			clientReceive.MessageType = content.MessageType
			clientReceive.Message = content.Message
			clientReceive.From.UserId = content.UserId
			clientReceive.From.UserType = content.UserType
			clientMessage, err := json.Marshal(clientReceive)
			if err !=nil{
				fmt.Println("rabbitmq消息转String出错",err)
			}
			userList := make([]string,0)
			if content.GroupSend == true{
				group := content.GroupTypeId + "_" +content.GroupId
				groupUser,err :=  RedisCache.SMembers(group).Result()
				if err!=nil{
					fmt.Println("redis缓存出错",err)
				}
				userList = append(userList, groupUser...)
			}
			userList = append(userList, content.SingleSend...)
			if userList !=nil{
				for _,uuid:= range userList {
					if ws, ok := models.ConnectMap[uuid]; ok {
						ws.WriteMessage(websocket.TextMessage, clientMessage)
						fmt.Println(uuid, "发送成功")
					}
				}
			}
			msg.Ack(false)
		}
	}()
}
func (r *Rabbitmq) SendMessage(queue_name,message string)  {
	err := r.Ch.Publish(
		"",     // 交换机名称
		queue_name, // 路由名称
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain", // 文本格式
			Body:        []byte(message), // 内容
		})

	if err!=nil{
		fmt.Println("RQ接受消息出错",err)
	}

}

func init()  {

	RQ.ConnectUrl=beego.AppConfig.String("cache.conn")
	RQ.Connect()
	RQ.Receive("task_queue")

}
