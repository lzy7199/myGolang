package CHandle

import (
	"Model"
	"Utils"
	"encoding/base64"
	"encoding/json"
	// "fmt"
	"io"
	"io/ioutil"
	// "log"
	// "fmt"
	"net"
	"net/http"
	"strings"
)

func Call(address string, method string, id interface{}, params interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	if err != nil {
		Utils.LogInfo("Marshal: %v", err)
		return nil, err
	}
	resp, err := http.Post(address, "application/json", strings.NewReader(string(data)))
	if err != nil {
		Utils.LogInfo("Post: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Utils.LogInfo("ReadAll: %v", err)
		return nil, err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		Utils.LogInfo("Unmarshal: %v", err)
		return nil, err
	}
	//log.Println(result)
	return result, nil
}

func CallTcp(conn *net.Conn, method string, uid interface{}, params interface{}) error {
	msg, err := json.Marshal(map[string]interface{}{
		"method": method,
		"uid":    uid,
		"params": params,
	})
	if err != nil {
		Utils.LogInfo("Marshal: %v", err)
		return err
	}
	// Utils.LogInfo("发送包json压缩成功！")

	// 发送给服务器
	errCode := writeToServer(conn, string(msg))
	if errCode != -1 {
		Utils.LogInfo("writeToServer error! errCode = %d\n", errCode)
		return Utils.LogErr(errCode)
	}
	Utils.LogInfo("msg = %s, 成功发送，等待回应！", msg)
	return nil
}

func CallTcp2(conn *net.Conn, method string, uid interface{}, params interface{}) error {
	_, err := json.Marshal(map[string]interface{}{
		"method": method,
		"uid":    uid,
		"params": params,
	})
	if err != nil {
		Utils.LogInfo("Marshal: %v", err)
		return err
	}
	// Utils.LogInfo("发送包json压缩成功！")

	// 发送给服务器
	// errCode := writeToServer(conn, string(msg))
	// if errCode != -1 {
	//     Utils.LogInfo("writeToServer error! errCode = %d\n", errCode)
	//     return Utils.LogErr(errCode)
	// }
	// Utils.LogInfo("msg = %s, 成功发送给服务器，等待回应！", msg)
	return nil
}

func HandleResult(conn *net.Conn, uid int64) {
	for {
		// 等待服务器返回
		length, err := readInt(conn)
		if err != nil {
			Utils.LogInfo("readInt error: %v", err)
			if err == io.EOF {
				Utils.LogInfo("服务器断开连接！")
				(*conn).Close()
			}
			break
		}
		Utils.LogInfo("chat result length = %d\n", length)

		// 读出加密过的数据包
		oriMsg, err := readStr(conn, length)
		if err != nil {
			Utils.LogInfo("readStr error: %v", err)
			if err == io.EOF {
				Utils.LogInfo("服务器断开连接！")
				(*conn).Close()
			}
			break
		}

		// base64解密
		finalMsg := make([]byte, 1024)
		_, err = base64.StdEncoding.Decode(finalMsg, []byte(oriMsg))
		if err != nil {
			Utils.LogInfo("base64 decode error! %s", err)
			break
		}
		finalMsgStr := strings.Trim(string(finalMsg), "\x00")
		// Utils.LogInfo("收到的包解密后为：%s\n", finalMsgStr)

		// json解密
		requestMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(finalMsgStr), &requestMap)
		if err != nil {
			Utils.LogInfo("json unmarshal error! %s", err)
			break
		}
		// Utils.LogInfo("收到的包json解密后为：%s\n", requestMap)

		// 根据收到的包进行不同的结果处理
		if methodInter, ok := requestMap["method"]; ok {
			switch methodInter.(string) {
			case "LoginChat":
				// 登陆成功，准备接收聊天信息
			case "SendChat":
				// 发送聊天是否成功
				errCode := int(requestMap["error"].(float64))
				if errCode != -1 {
					Utils.LogInfo("返回值出错啦！err = %s", Utils.LogErr(errCode).Error())
				}
			case "ReceiveChat":
				// 收到聊天信息
				resultMsgMap := requestMap["result"].(map[string]interface{})
				chatMsg := Model.ChangeMapToChatMsg(resultMsgMap)
				switch chatMsg.ChatType {
				case 1:
					Utils.LogInfo("“%s对大家说：%s”", chatMsg.FromName, chatMsg.Msg)
				case 2:
					Utils.LogInfo("“%s对你说：%s”", chatMsg.FromName, chatMsg.Msg)
				case 3:
					Utils.LogInfo("“系统公告：%s”", chatMsg.Msg)
				}
			case "IsHeartBeat":
				Utils.LogInfo("收到心跳包！")
				// 发送心跳包回应
				err = CallTcp(conn, "IsHeartBeat", uid, make(map[string]interface{}))
				if err != nil {
					Utils.LogInfo("CallTcp（IsHeartBeat）出错啦！err = %s", err.Error())
				}
			default:
				// 未知方法
				Utils.LogInfo("”%s“是什么方法？", methodInter.(string))
			}
		} else {
			Utils.LogInfo("你妹啊！都没方法的！")
		}
	}
}

/**
向聊天服务器发送消息
param： conn 与聊天服务器连接的客户端
        msg 欲发送的消息包
result：int 错误编号
**/
func writeToServer(conn *net.Conn, msg string) int {
	// base64加密
	result := base64.StdEncoding.EncodeToString([]byte(msg))

	// 加上长度，并发送
	dataLenByte := Utils.Uint32ToBytes(uint32(len(result)))
	Utils.LogInfo("发送的包长度为%d\n", len(result))
	sendMsg := append(dataLenByte[:], result...)
	_, err := (*conn).Write(sendMsg)
	if err != nil {
		Utils.LogErr(err)
		if err == io.EOF {
			Utils.LogInfo("服务器断开连接！")
			(*conn).Close()
		}
		return 110
	}
	return -1
}

/**
从conn中读取一个int32值
param：conn 客户端连接socket对象
result：int 读取的int32值
        error 错误对象
**/
func readInt(conn *net.Conn) (int, error) {
	data, err := Utils.ReadConn(conn, 4)
	if err != nil {
		return 0, err
	}
	//  Utils.LogInfo("read data=%#v\n",data)
	return int(Utils.BytesToUint32(data)), nil
}

/**
从conn中读取指定长度的字符串
param：conn 客户端连接socket对象
        num 指定长度
result：string 读取的字符串
        error 错误对象
**/
func readStr(conn *net.Conn, num int) (string, error) {
	data, err := Utils.ReadConn(conn, num)
	if err != nil {
		return "", err
	}
	//  Utils.LogInfo("read datastr=%v\n",data)
	// return strings.Trim(string(data), "\r\n\t "), nil
	return string(data), nil
}
