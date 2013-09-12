package main

import (
	"CHandle"
	"Config"
	"Model"
	"Utils"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type User struct {
	Id       int64
	Username string
	Password string
}

func main() {
	// 定义日志文件，激活配置文件
	Utils.LogInfo("init log！\n")
	Utils.LogInfo("log file=%s\n", Config.GetLogOutFile())
	err := Utils.InitLogOut(Config.GetLogOutFile())
	if err != nil {
		return
	}
	defer Utils.DeferFiles()

	//创建守护进程
	Utils.LogInfo("client main start！", os.Getppid())
	Utils.LogInfo("ppid=%d\n", os.Getppid())
	if os.Getppid() != 1 {
		filePath, _ := filepath.Abs(os.Args[0])
		Utils.LogInfo("filePath=%d\n", filePath)
		cmd := exec.Command(filePath, os.Args[1:]...)
		cmd.Start()
		return
	}

	var add int64 = 0
	var i int64 = 1 + add

	// 测试http
	// for {
	// 	go doFun(i)
	// 	i = i + 100000
	// 	time.Sleep(1e7)
	// }
	// 测试tcp
	for {
		if i <= 50+add {
			go doTcpTest(i)
			i++
		}
		time.Sleep(1e9)
	}
}

func doTcpTest(uid int64) {
	// conn, err := net.Dial("tcp", "127.0.0.1:12321")
	conn, err := net.Dial("tcp", "192.168.0.220:12321")
	if err != nil {
		Utils.LogInfo("出错啦！err = %s", err.Error())
		return
	}
	Utils.LogInfo("net.Dial成功！")
	// 登陆
	chatAvatar := new(Model.ChatAvatar)
	chatAvatar.AvatarId = uid
	chatAvatar.Name = fmt.Sprintf("无敌小宇%d号", uid)
	chatAvatar.Uid = uid
	chatAvatar.ServerIndex = 0
	err = CHandle.CallTcp(&conn, "LoginChat", uid, Model.ChangeChatAvatarToMap(chatAvatar))
	if err != nil {
		Utils.LogInfo("出错啦！err = %s", err.Error())
	}
	// 开启收取线程
	go CHandle.HandleResult(&conn, uid)

	// 循环发送世界聊天
	// 开始世界聊天
	i := 1
	chatMsg := new(Model.ChatMsg)
	for {
		chatMsg.FromName = fmt.Sprintf("无敌小宇%d号", uid)
		chatMsg.ChatType = 1
		chatMsg.Msg = fmt.Sprintf("我第%d次说话了！！大家看得见吗？", i)
		err = CHandle.CallTcp(&conn, "SendChat", uid, Model.ChangeChatMsgToMap(chatMsg))
		if err != nil {
			Utils.LogInfo("CallTcp出错啦！err = %s", err.Error())
			break
		}
		i++
		time.Sleep(2e8)
	}
}

func doFun(j int64) {
	for i := j; i < j+100000; i++ {
		user := new(User)
		user.Id = 0
		user.Username = fmt.Sprintf("lzy%d", i)
		user.Password = strconv.Itoa(rand.Intn(10000000))
		// res, err:=Handle.Call("http://user:pass@127.0.0.1:8332", "getinfo", 1, []interface{}{})
		res, err := CHandle.Call("http://root:srl001988@192.168.0.220:12345", "RegisterUser", i, user)
		// res, err := CHandle.Call("http://127.0.0.1:12345", "RegisterUser", i, user)
		if err != nil {
			Utils.LogInfo("Err: %v", err)
		}
		log.Println(res)
	}
}
