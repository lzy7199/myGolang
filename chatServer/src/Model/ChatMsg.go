package Model

import (
	"Utils"
	"strconv"
)

type ChatMsg struct {
	ChatType int
	FromName string
	ToName   string
	Msg      string
}

func ChangeMapToChatMsg(m map[string]interface{}) *ChatMsg {
	chatMsg := new(ChatMsg)
	chatMsg.ChatType = Utils.ChangeMapToPro(m["chatType"], "int").(int)
	chatMsg.FromName = Utils.ChangeMapToPro(m["fromName"], "string").(string)
	chatMsg.ToName = Utils.ChangeMapToPro(m["toName"], "string").(string)
	chatMsg.Msg = Utils.ChangeMapToPro(m["msg"], "string").(string)
	return chatMsg
}

func ChangeChatMsgToMap(chatMsg *ChatMsg) map[string]string {
	m := make(map[string]string)
	m["chatType"] = strconv.Itoa(chatMsg.ChatType)
	m["fromName"] = chatMsg.FromName
	m["toName"] = chatMsg.ToName
	m["msg"] = chatMsg.Msg
	return m
}
