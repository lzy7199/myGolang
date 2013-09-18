package Utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"runtime"
	"strconv"
)

/**
用gob进行数据编码
**/
func GobEncode(data interface{}) ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(data)
	if err != nil {
		LogErr(err)
		return nil, LogErr(103)
	}
	return network.Bytes(), nil
}

/**
用gob进行数据解码
**/
func GobDecode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(to)
	if err != nil {
		LogErr(err)
		return LogErr(104)
	}
	return err
}

/**
从int32转化为[]byte
**/
func Uint32ToBytes(i uint32) []byte {
	return []byte{byte((i >> 24) & 0xff), byte((i >> 16) & 0xff),
		byte((i >> 8) & 0xff), byte(i & 0xff)}
}

/**
从[]byte转化为int32
**/
func BytesToUint32(buf []byte) uint32 {
	return uint32(buf[0])<<24 + uint32(buf[1])<<16 + uint32(buf[2])<<8 +
		uint32(buf[3])
}

/**
用json进行数据解码
**/
func JsonDecode(byteArr []byte) (map[string]interface{}, error) {
	var msg interface{}
	err := json.Unmarshal(byteArr, &msg)
	if err != nil {
		return nil, LogErr(err)
	}
	return msg.(map[string]interface{}), nil
}

/**
用json进行数据编码
**/
func JsonEncode(jsonData interface{}) ([]byte, error) {
	msg, err := json.Marshal(jsonData)
	if err != nil {
		LogErr(102)
		return nil, err
	}
	return msg, err
}

/**
从net.TCPConn 中读取固定长度（防止长短包粘包）
**/
func ReadConn(conn *net.TCPConn, readLen int) ([]byte, error) {
	//	LogInfo("need read data=%d\n",readLen)
	dataBuf := make([]byte, readLen)
	var dataLenTag int
	for {
		tmpNum, err := (*conn).Read(dataBuf[dataLenTag:readLen])
		if err != nil {
			if err == io.EOF {
				LogInfo("read EOF  num=%d\n", tmpNum)
				return dataBuf, err
			}
			if err == io.ErrUnexpectedEOF {
				LogInfo("read ErrUnexpectedEOF  num=%d\n", tmpNum)
				return dataBuf, err
			}
			//			LogInfo("read num=%d\n", tmpNum)
			LogInfo("err info=%v\n", err)
			if err.Error() == "use of closed network connection" {
				err = io.ErrClosedPipe
			}
			return dataBuf, err
		}
		// LogInfo("read num=%d\n", tmpNum)
		// LogInfo("read data=%v\n", dataBuf)
		dataLenTag += tmpNum
		if dataLenTag >= readLen {
			break
		}
	}
	return dataBuf, nil
}

/**
将map["XXX"]转化为typeStr所表示的格式值
param：src 源值
		typeStr 转化格式
result：转化后的值
**/
func ChangeMapToPro(src interface{}, typeStr string) interface{} {
	def := false
	if src == nil {
		def = true
	}
	switch typeStr {
	case "string":
		if def {
			return ""
		} else {
			return src.(string)
		}
	case "int", "int32":
		if def {
			return 0
		} else {
			result, err := strconv.Atoi(src.(string))
			if err != nil {
				LogErr(err)
			}
			return result
		}
	case "int64":
		if def {
			return int64(0)
		} else {
			result, err := strconv.ParseInt(src.(string), 10, 64)
			if err != nil {
				LogErr(err)
			}
			return result
		}
	case "float", "float32":
		if def {
			return float32(0)
		} else {
			result, err := strconv.ParseFloat(src.(string), 32)
			if err != nil {
				LogErr(err)
			}
			return result
		}
	case "float64":
		if def {
			return float64(0)
		} else {
			result, err := strconv.ParseFloat(src.(string), 64)
			if err != nil {
				LogErr(err)
			}
			return result
		}
	case "[]int64":
		return src.([]int64)
	case "[]int", "[]int32":
		return src.([]int)
	}
	return nil
}

func CheckPanic(finalErr interface{}) {
	var stack string
	LogErrInfo("------------------handler crashed with error (%s)------------------", finalErr)
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// LogErrInfo("------------------file（%s）line(%d) error!", file, line)
		stack = stack + fmt.Sprintln(file, line)
	}
	LogErrInfo("------------------------------------")
	LogErrInfo("stack = (%s)", stack)
	LogErrInfo("------------------------------------")
}
