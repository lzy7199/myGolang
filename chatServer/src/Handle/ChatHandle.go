package Handle

import (
	"Utils"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	// "strconv"
	"Model"
	// "fmt"
	"strings"
	"time"
)

/**聊天同步模型类对象**/
// var chatChan chan *Model.ChatSyncModel
var chatSyncModelTotal *Model.ChatSyncModel

/**多路复用，用于跟踪每一个调用特定的回调函数**/
type ServeMux struct {
	m map[string]func(*net.TCPConn, map[string]interface{}) map[string]interface{}
}

func init() {
	mainMux.m = make(map[string]func(*net.TCPConn, map[string]interface{}) map[string]interface{})
	// chatChan = make(chan *Model.ChatSyncModel, 1)
	chatSyncModelTotal = new(Model.ChatSyncModel)
	chatSyncModelTotal.SocketMap = make(map[*net.TCPConn](*Model.ChatAvatar))
	chatSyncModelTotal.NameMap = make(map[int]map[string](*net.TCPConn))
	// chatChan <- chatSyncModel
	go tickHeartBeat()
}

/**
初始化心跳包机制
param：conn 客户端连接socket对象
result：int 读取的int32值
		error 错误对象
**/
func tickHeartBeat() {
	defer func() {
		if finalErr := recover(); finalErr != nil {
			Utils.LogErrInfo("tickHeartBeat finalErr = %s", finalErr)
			Utils.CheckPanic(finalErr)
		}
	}()
	c := time.Tick(5 * time.Minute)
	for _ = range c {
		// 定时循环发送心跳包
		checkHeartBeat()
	}
}

func checkHeartBeat() {
	chatSyncModelTotal.M.Lock()
	defer chatSyncModelTotal.M.Unlock()
	Utils.LogInfo("发心跳包，当前登陆人数：%d", len(chatSyncModelTotal.SocketMap))
	for conn, ca := range chatSyncModelTotal.SocketMap {
		// 判断是发送心跳包还是断开连接
		if ca.IsAlive == 1 {
			// 发送心跳包
			ca.IsAlive = 0
			writeBackSuccessWithoutLock(conn, createResult("IsHeartBeat", -1, nil, -1), chatSyncModelTotal)
		} else {
			// 没反应，直接断开
			// exitMapWithoutLock(conn, chatSyncModelTotal)
			Utils.LogErrInfo("%d号童鞋没发心跳包给服务器哦！", ca.Uid)
		}
	}
}

// for {
// 	// 定时循环发送心跳包
// 	func() {
// 		chatSyncModel := getChatSyncModel()
// 		defer setChatSyncModel(chatSyncModel)
// 		Utils.LogInfo("发心跳包，当前登陆人数：%d", len(chatSyncModel.SocketMap))
// 		for conn, ca := range chatSyncModel.SocketMap {
// 			// 判断是发送心跳包还是断开连接
// 			if ca.IsAlive == 1 {
// 				// 发送心跳包
// 				ca.IsAlive = 0
// 				writeBackSuccess(conn, createResult("IsHeartBeat", -1, nil, -1))
// 			} else {
// 				// 没反应，直接断开
// 				exitMapWithoutLock(conn, chatSyncModel)
// 			}
// 		}
// 	}()
// 	time.Sleep(60e9)
// }

// func pushChatChan(chatSyncModel *Model.ChatSyncModel) {
// 	Utils.LogInfo("setChatSyncModel, num = %d", len(chatChan))
// 	if len(chatChan) == 0 {
// 		chatChan <- chatSyncModel
// 	}
// }

// func popChatChan() *Model.ChatSyncModel {
// 	Utils.LogInfo("getChatSyncModel, num = %d", len(chatChan))
// 	return <-chatChan
// }

// func getChatSyncModel() *Model.ChatSyncModel {
// 	return chatSyncModelTotal
// }

func GetCurClientNum() int {
	// defer func() {
	// 	if finalErr := recover(); finalErr != nil {
	// 		Utils.LogErrInfo("GetCurClientNum finalErr = %s", finalErr)
	// 		Utils.CheckPanic(finalErr)
	// 	}
	// }()
	// Utils.LogInfo("chatSyncModelTotal = %s", chatSyncModelTotal)
	// defer setChatSyncModel(chatSyncModel)
	return len(chatSyncModelTotal.SocketMap)
}

/**
注册所有监听函数
**/
func RegisterAllFunc() {
	// --------------登陆管理模块--------------
	handleFunc("LoginChat", LoginChat)
	handleFunc("LogoutChat", LogoutChat)

	// --------------聊天管理模块--------------
	handleFunc("SendChat", SendChat)

	// --------------其他模块--------------
	handleFunc("IsHeartBeat", IsHeartBeat)
	handleFunc("GetClientNum", GetClientNum)
}

/**
多路复用multiplexer的实例
**/
var mainMux ServeMux

/**
用于注册回调函数
**/
func handleFunc(pattern string, handler func(*net.TCPConn, map[string]interface{}) map[string]interface{}) {
	mainMux.m[pattern] = handler
}

/**
处理客户端发来的聊天命令
param：conn 客户端连接socket对象
		clientIndex 客户端编号
**/
func HandleChat(conn *net.TCPConn, clientIndex int) {
	// bytes, err := Utils.ReadConn(conn, 300)
	// if err != nil {
	// 	Utils.LogErr(err)
	// 	fmt.Printf("bytes = %s", string(bytes))
	// 	Utils.LogInfo("bytes = %s", string(bytes))
	// } else {
	// 	fmt.Printf("bytes = %s", string(bytes))
	// 	Utils.LogInfo("bytes = %s", string(bytes))
	// }
	defer func() {
		if finalErr := recover(); finalErr != nil {
			Utils.LogErrInfo("HandleChat finalErr = %s", finalErr)
			Utils.CheckPanic(finalErr)
		}
	}()
	for {
		// 读出长度
		length, err := readInt(conn)
		if err != nil {
			Utils.LogInfo("readInt error, error = %s", err.Error())
			writeBackException("", conn, 20001, -1)
			exitMap(conn)
			break
		}
		// Utils.LogInfo("收到的包长度为：%d\n", length)
		if length > 10000000 {
			Utils.LogErrCode(20005)
			writeBackException("", conn, 20005, -1)
			exitMap(conn)
			break
		}
		// 读出包内容
		chatStr, err := readStr(conn, length)
		if err != nil {
			Utils.LogInfo("readInt error, error = %s", err.Error())
			writeBackException("", conn, 20001, -1)
			continue
		}
		// Utils.LogInfo("收到的包解密前为：%s\n", chatStr)

		// base64解码
		finalChat := make([]byte, 1024)
		_, err = base64.StdEncoding.Decode(finalChat, []byte(chatStr))
		if err != nil {
			Utils.LogInfo("base64 decode error! %s", err)
			writeBackException("", conn, 20002, -1)
			continue
		}
		finalChatStr := strings.Trim(string(finalChat), "\x00")
		// Utils.LogInfo("收到的包解密后为：%s\n", finalChatStr)

		// json解码
		requestMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(finalChatStr), &requestMap)
		if err != nil {
			Utils.LogInfo("json Unmarshal error, error = %s", err.Error())
			writeBackException("", conn, 20003, -1)
			continue
		}
		// Utils.LogInfo("收到的包json解密后为：%s\n", requestMap)

		// 判断所调用的方法函数
		if _, ok := requestMap["method"]; !ok {
			writeBackException("", conn, 100, requestMap["uid"])
		}
		function, ok := mainMux.m[requestMap["method"].(string)]
		if ok {
			//如果函数存在，则调用它
			response := function(conn, requestMap)
			if response != nil {
				if response["error"].(int) == -1 {
					// 成功
					writeBackSuccess(conn, response)
				} else {
					// 出错
					writeBackException(requestMap["method"].(string), conn, response["error"].(int), requestMap["uid"])
				}
			}
		} else {
			//如果函数不存在
			Utils.LogInfo("HTTP JSON RPC Handle - No function to call for ‘%s’", requestMap["method"])
			writeBackException(requestMap["method"].(string), conn, 100, requestMap["uid"])
		}
	}
	Utils.LogInfo("clinet--%d is disConnected！now total client = %d\n", clientIndex, GetCurClientNum())
}

/**
从conn中读取一个int32值
param：conn 客户端连接socket对象
result：int 读取的int32值
		error 错误对象
**/
func readInt(conn *net.TCPConn) (int, error) {
	data, err := Utils.ReadConn(conn, 4)
	if err != nil {
		return 0, err
	}
	//	Utils.LogInfo("read data=%#v\n",data)
	return int(Utils.BytesToUint32(data)), nil
}

/**
从conn中读取指定长度的字符串
param：conn 客户端连接socket对象
		num 指定长度
result：string 读取的字符串
		error 错误对象
**/
func readStr(conn *net.TCPConn, num int) (string, error) {
	data, err := Utils.ReadConn(conn, num)
	if err != nil {
		return "", err
	}
	//	Utils.LogInfo("read datastr=%v\n",data)
	// return strings.Trim(string(data), "\r\n\t "), nil
	return string(data), nil
}

/**
将数据通过base64加密后发送回客户端
**/
func writeBack(conn *net.TCPConn, data []byte) {
	// base64加密字符串
	result := base64.StdEncoding.EncodeToString(data)
	// 加长度
	dataLenByte := Utils.Uint32ToBytes(uint32(len(result)))
	// Utils.LogInfo("packageLen：%d\n", len(result))
	data = append(dataLenByte[:], []byte(result)...)

	_, err := (*conn).Write(data)
	if err != nil {
		Utils.LogInfo("writeBack error, error = %s", err.Error())
		checkIsExitMap(conn, err)
	}
}

/**
将数据通过base64加密后发送回客户端
**/
func writeBackWithoutLock(conn *net.TCPConn, data []byte, chatSyncModel *Model.ChatSyncModel) {
	// base64加密字符串
	result := base64.StdEncoding.EncodeToString(data)
	// 加长度
	dataLenByte := Utils.Uint32ToBytes(uint32(len(result)))
	// Utils.LogInfo("packageLen：%d\n", len(result))
	data = append(dataLenByte[:], []byte(result)...)

	_, err := (*conn).Write(data)
	if err != nil {
		Utils.LogInfo("writeBackWithoutLock error, error = %s", err.Error())
		checkIsexitMapWithoutLock(conn, err, chatSyncModel)
	}
}

/**
将错误编号通过base64加密后发送回客户端
**/
func writeBackException(method interface{}, conn *net.TCPConn, errCode int, uid interface{}) {
	data, err := json.Marshal(createResult(method, nil, Utils.LogErrCode(errCode), uid))
	if err != nil {
		Utils.LogInfo("writeBackException error, error = %s", err.Error())
		Utils.LogErrCode(107)
		return
	}
	writeBack(conn, data)
}

/**
将错误编号通过base64加密后发送回客户端
**/
func writeBackExceptionWithoutLock(method interface{}, conn *net.TCPConn, errCode int, uid interface{}, chatSyncModel *Model.ChatSyncModel) {
	data, err := json.Marshal(createResult(method, nil, Utils.LogErrCode(errCode), uid))
	if err != nil {
		Utils.LogInfo("writeBackExceptionWithoutLock error, error = %s", err.Error())
		Utils.LogErrCode(107)
		return
	}
	writeBackWithoutLock(conn, data, chatSyncModel)
}

/**
将成功返回值通过base64加密后发送回客户端
**/
func writeBackSuccess(conn *net.TCPConn, response map[string]interface{}) {
	data, err := json.Marshal(response)
	if err != nil {
		Utils.LogInfo("writeBackSuccess error, error = %s", err.Error())
		Utils.LogErrCode(107)
		return
	}
	writeBack(conn, data)
}

/**
将成功返回值通过base64加密后发送回客户端
**/
func writeBackSuccessWithoutLock(conn *net.TCPConn, response map[string]interface{}, chatSyncModel *Model.ChatSyncModel) {
	data, err := json.Marshal(response)
	if err != nil {
		Utils.LogInfo("writeBackSuccessWithoutLock error, error = %s", err.Error())
		Utils.LogErrCode(107)
		return
	}
	writeBackWithoutLock(conn, data, chatSyncModel)
}

/**
创建返回结果
param： result 欲返回的结果
        err 欲返回的错误编号
        uid 当前用户的uid
result：返回结果
**/
func createResult(method, result, err, uid interface{}) map[string]interface{} {
	if err == nil {
		err = -1
	}
	return map[string]interface{}{
		"method": method,
		"result": result,
		"error":  err,
		"uid":    uid,
	}
}

func checkIsExitMap(conn *net.TCPConn, err error) {
	if err == io.EOF || err == io.ErrClosedPipe || strings.Contains(err.Error(), "broken pipe") {
		exitMap(conn)
	}
}

func checkIsexitMapWithoutLock(conn *net.TCPConn, err error, chatSyncModel *Model.ChatSyncModel) {
	if err == io.EOF || err == io.ErrClosedPipe || strings.Contains(err.Error(), "pipe") || strings.Contains(err.Error(), "use of closed network connection") {
		exitMapWithoutLock(conn, chatSyncModel)
	}
}
