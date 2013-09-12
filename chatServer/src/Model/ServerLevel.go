package Model

import (
	"Utils"
	"strconv"
)

type ServerLevel struct {
	ServerAddress string // 服务器地址
	ServerId      int    // 服务器编号
	ServerName    string // 服务器名称
	CurStatus     int    // 0:停机维护（维护）（维护时所有字体变成灰色）1：正常（不写字）2：推荐服务器（荐）3：拥挤的服务器（拥挤）4：活动服务器（活动）5：测试服务器（测试）6：新服务器（新）
}

func ChangeMapToServerLevel(m map[string]interface{}) *ServerLevel {
	serverLevel := new(ServerLevel)
	serverLevel.ServerAddress = Utils.ChangeMapToPro(m["serverAddress"], "string").(string)
	serverLevel.CurStatus = Utils.ChangeMapToPro(m["curStatus"], "int").(int)
	serverLevel.ServerId = Utils.ChangeMapToPro(m["serverId"], "int").(int)
	serverLevel.ServerName = Utils.ChangeMapToPro(m["serverName"], "string").(string)
	return serverLevel
}

func ChangeServerLevelToMap(serverLevel *ServerLevel) map[string]string {
	m := make(map[string]string)
	m["curStatus"] = strconv.Itoa(serverLevel.CurStatus)
	m["serverAddress"] = serverLevel.ServerAddress
	m["serverName"] = serverLevel.ServerName
	m["serverId"] = strconv.Itoa(serverLevel.ServerId)
	return m
}
