package Handle

import (
	// "DataStore"
	// "Lib/EmailSend"
	"Model"
	"Utils"
	"net"
	// "strconv"
)

/**
发送聊天信息
param： conn socket连接对象
        m 请求的参数
result：返回结果
**/
func SendChat(conn *net.TCPConn, m map[string]interface{}) map[string]interface{} {
	if params, ok := m["params"]; ok {
		chatMsg := Model.ChangeMapToChatMsg(params.(map[string]interface{}))
		if chatSyncModelTotal.SocketMap[conn] == nil {
			return createResult("SendChat", -1, nil, m["uid"])
		}
		switch chatMsg.ChatType {
		case 1:
			// Utils.LogInfo("发世界聊天啦！-------------")
			// 本服世界聊天
			response := createResult("ReceiveChat", Model.ChangeChatMsgToMap(chatMsg), nil, m["uid"])
			localNameMap := chatSyncModelTotal.NameMap[chatSyncModelTotal.SocketMap[conn].ServerIndex]
			// Utils.LogInfo("localNameMap = %s", localNameMap)
			for _, connTemp := range localNameMap {
				// 循环发送
				// Utils.LogInfo("conn = %s, connTemp = %s", conn, connTemp)
				writeBackSuccess(connTemp, response)
			}
		case 2:
			// 本服私聊
			// 判断对方是否存在
			if tarConn, ok := chatSyncModelTotal.NameMap[chatSyncModelTotal.SocketMap[conn].ServerIndex][chatMsg.ToName]; ok {
				writeBackSuccess(tarConn, createResult("ReceiveChat", Model.ChangeChatMsgToMap(chatMsg), nil, m["uid"]))
			} else {
				return createResult("SendChat", nil, Utils.LogErrCode(20007), m["uid"])
			}
		case 3:
			// 本服系统公告
			response := createResult("ReceiveChat", Model.ChangeChatMsgToMap(chatMsg), nil, m["uid"])
			localNameMap := chatSyncModelTotal.NameMap[chatSyncModelTotal.SocketMap[conn].ServerIndex]
			for _, conn := range localNameMap {
				// 循环发送
				writeBackSuccess(conn, response)
			}
		case 4:
			// 所有服聊天
			response := createResult("ReceiveChat", Model.ChangeChatMsgToMap(chatMsg), nil, m["uid"])
			for conn, _ := range chatSyncModelTotal.SocketMap {
				// 循环发送
				writeBackSuccess(conn, response)
			}
		case 5:
			// 所有服系统公告
			response := createResult("ReceiveChat", Model.ChangeChatMsgToMap(chatMsg), nil, m["uid"])
			for conn, _ := range chatSyncModelTotal.SocketMap {
				// 循环发送
				writeBackSuccess(conn, response)
			}
		default:
			return createResult("SendChat", nil, Utils.LogErrCode(20004), m["uid"])
		}
		// 打印日志
		Utils.LogInfo("SendChat--------connIp:%s, params:%s", (*conn).RemoteAddr(), params)
		return createResult("SendChat", -1, nil, m["uid"])
	} else {
		return createResult("SendChat", nil, Utils.LogErrCode(101), m["uid"])
	}
	return nil
}
