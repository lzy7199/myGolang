package Utils

import (
	"fmt"
)

var err map[int]string

type ServerErr struct {
	ErrCode int
	ErrMsg  string
}

func (e ServerErr) Error() string {
	return fmt.Sprintf("%v: %v", e.ErrCode, e.ErrMsg)
}

func init() {
	err = map[int]string{
		/*******************chat服务器错误*****************/
		20001: "chat read error",                 //聊天语句读出异常（流异常）
		20002: "base64 decode error",             //base64解码异常
		20003: "json decode error",               //json解码异常
		20004: "chatType error",                  //无此聊天类型
		20005: "package is too long error",       // 包太长
		20006: "repeat login error",              //重复登陆，踢除
		20007: "target is not exist error",       //对方不存在或未在线
		20008: "client already over limit error", //该服务器登陆人数超过上限，请稍后再试
	}
}

/**
打印错误
**/
func LogErr(errCode interface{}) error {
	if errCode != nil {
		switch value := errCode.(type) {
		case int:
			//			value:=errCode.(errCode)
			errMsg, exist := err[value]
			if exist == true {
				LogInfo("err=%s\n", fmt.Sprintf("errCode=%d,errMsg=%s", value, errMsg))
				return ServerErr{value, errMsg}
			} else {
				LogInfo("err=%s\n", "unknow err")
				return ServerErr{1000, "unknow err"}
			}
		case error:
			LogInfo("err=%s\n", errCode)
			return errCode.(error)
		}
	}
	return nil
}

/**
根据错误编号打印错误
**/
func LogErrCode(errCode int) int {
	errMsg, exist := err[errCode]
	if exist == true {
		LogInfo("err=%s\n", fmt.Sprintf("errCode=%d,errMsg=%s", errCode, errMsg))
	} else {
		LogInfo("err=%s,%d\n", "unknow err", errCode)
	}
	return errCode
}

/**
打印并抛出错误
**/
func LogPanicErr(err interface{}) {
	panic(err)
}
