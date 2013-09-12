package Model

import (
	"Utils"
	"strconv"
)

type ChatAvatar struct {
	AvatarId    int64
	Uid         int64
	Name        string
	ServerIndex int
	IsAlive     int
}

func ChangeMapToChatAvatar(m map[string]interface{}) *ChatAvatar {
	chatAvatar := new(ChatAvatar)
	chatAvatar.AvatarId = Utils.ChangeMapToPro(m["avatarId"], "int64").(int64)
	chatAvatar.Uid = Utils.ChangeMapToPro(m["uid"], "int64").(int64)
	chatAvatar.Name = Utils.ChangeMapToPro(m["name"], "string").(string)
	chatAvatar.ServerIndex = Utils.ChangeMapToPro(m["serverIndex"], "int").(int)
	chatAvatar.IsAlive = Utils.ChangeMapToPro(m["isAlive"], "int").(int)
	return chatAvatar
}

func ChangeChatAvatarToMap(chatAvatar *ChatAvatar) map[string]string {
	m := make(map[string]string)
	m["avatarId"] = strconv.FormatInt(chatAvatar.AvatarId, 10)
	m["uid"] = strconv.FormatInt(chatAvatar.Uid, 10)
	m["name"] = chatAvatar.Name
	m["serverIndex"] = strconv.Itoa(chatAvatar.ServerIndex)
	m["isAlive"] = strconv.Itoa(chatAvatar.IsAlive)
	return m
}
