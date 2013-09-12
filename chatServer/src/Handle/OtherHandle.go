package Handle

import (
	"net"
)

/**
正在心跳
param： conn socket连接对象
        m 请求的参数
result：返回结果
**/
func IsHeartBeat(conn *net.TCPConn, m map[string]interface{}) map[string]interface{} {
	chatSyncModel := popChatChan()
	defer pushChatChan(chatSyncModel)
	if ca, ok := chatSyncModel.SocketMap[conn]; ok {
		ca.IsAlive = 1
	}
	return nil
}
