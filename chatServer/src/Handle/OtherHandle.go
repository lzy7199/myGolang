package Handle

import (
	"Utils"
	"net"
)

/**
正在心跳
param： conn socket连接对象
        m 请求的参数
result：返回结果
**/
func IsHeartBeat(conn *net.TCPConn, m map[string]interface{}) map[string]interface{} {
	chatSyncModelTotal.M.Lock()
	defer chatSyncModelTotal.M.Unlock()
	if ca, ok := chatSyncModelTotal.SocketMap[conn]; ok {
		ca.IsAlive = 1
	}
	return nil
}

/**
获得当前登陆的客户端数量
param： conn socket连接对象
        m 请求的参数
result：返回结果
**/
func GetClientNum(conn *net.TCPConn, m map[string]interface{}) map[string]interface{} {
	curClientNum := GetCurClientNum()
	// 打印日志
	Utils.LogInfo("GetClientNum--------totalClient:%s", curClientNum)
	return createResult("GetClientNum", curClientNum, nil, m["uid"])
}
