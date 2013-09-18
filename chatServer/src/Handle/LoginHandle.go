package Handle

import (
	// "DataStore"
	// "Lib/EmailSend"
	"Config"
	"Model"
	"Utils"
	"net"
	// "strconv"
)

/**
登陆聊天服
param： conn socket连接对象
        m 请求的参数
result：返回结果
**/
func LoginChat(conn *net.TCPConn, m map[string]interface{}) map[string]interface{} {
	if params, ok := m["params"]; ok {
		chatAvatar := Model.ChangeMapToChatAvatar(params.(map[string]interface{}))
		chatSyncModelTotal.M.Lock()
		defer chatSyncModelTotal.M.Unlock()
		// 判断登陆人数是否已经到达上限
		if len(chatSyncModelTotal.NameMap[chatAvatar.ServerIndex]) >= Config.GetLogicServerList()[chatAvatar.ServerIndex].MaxClient {
			return createResult("LoginChat", nil, Utils.LogErrCode(20008), m["uid"])
		}
		// 注册map
		chatAvatar.IsAlive = 1
		chatSyncModelTotal.SocketMap[conn] = chatAvatar
		nameMapMap, ok := chatSyncModelTotal.NameMap[chatAvatar.ServerIndex]
		if !ok {
			nameMapMap = make(map[string](*net.TCPConn))
		}
		if oldConn, ok := nameMapMap[chatAvatar.Name]; ok {
			// 踢出原来的玩家
			writeBackExceptionWithoutLock("LoginChat", oldConn, 20006, -1, chatSyncModelTotal)
		}
		nameMapMap[chatAvatar.Name] = conn
		chatSyncModelTotal.NameMap[chatAvatar.ServerIndex] = nameMapMap
		// 打印日志
		Utils.LogInfo("LoginChat--------connIp:%s, params:%s", (*conn).RemoteAddr(), params)
		return createResult("LoginChat", -1, nil, m["uid"])
	} else {
		return createResult("LoginChat", nil, Utils.LogErrCode(101), m["uid"])
	}
	return nil
}

/**
登出聊天服
param： conn socket连接对象
        m 请求的参数
result：返回结果
**/
func LogoutChat(conn *net.TCPConn, m map[string]interface{}) map[string]interface{} {
	// 退出map
	exitMap(conn)
	// 打印日志
	Utils.LogInfo("LogoutChat--------connIp:%s", (*conn).RemoteAddr())
	return createResult("LogoutChat", -1, nil, m["uid"])
}

/**
退出map（加锁）
param： conn socket连接对象
**/
func exitMap(conn *net.TCPConn) {
	chatSyncModelTotal.M.Lock()
	defer chatSyncModelTotal.M.Unlock()
	exitMapWithoutLock(conn, chatSyncModelTotal)
}

/**
退出map（不加锁）
param： conn socket连接对象
**/
func exitMapWithoutLock(conn *net.TCPConn, chatSyncModel *Model.ChatSyncModel) {
	// 退出map
	if _, ok := chatSyncModel.SocketMap[conn]; ok {
		name := chatSyncModel.SocketMap[conn].Name
		serverIndex := chatSyncModel.SocketMap[conn].ServerIndex
		delete(chatSyncModel.SocketMap, conn)
		delete(chatSyncModel.NameMap[serverIndex], name)
	}
	// 断开连接
	(*conn).Close()
}
