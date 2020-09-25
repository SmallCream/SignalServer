package models

import "github.com/gorilla/websocket"

var ConnectMap = make(map[string] *websocket.Conn)

type User struct {
	UserId      string
	UserType    string
	GroupTypeId string
	GroupId     string
}
type Content struct {
	MessageType string
	Message string
	UserId      string
	UserType    string
	GroupTypeId string
	GroupId     string
	SingleSend 	[]string
	GroupSend bool
}

type FromUser struct {
	UserId string
	UserType string
}

type ClientReceive struct {
	MessageType  string
	Message  string
	From    FromUser
}

