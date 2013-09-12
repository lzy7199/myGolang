package main

import (
	"Config"
	"Handle"
	"Utils"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

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
	Utils.LogInfo("chatserver main start！", os.Getppid())
	Utils.LogInfo("ppid=%d\n", os.Getppid())
	if os.Getppid() != 1 {
		filePath, _ := filepath.Abs(os.Args[0])
		Utils.LogInfo("filePath=%d\n", filePath)
		cmd := exec.Command(filePath, os.Args[1:]...)
		cmd.Start()
		return
	}

	// 注册监听函数
	Handle.RegisterAllFunc()

	// 执行信号量处理线程
	Utils.LogInfo("init signalHandle！\n")
	go signalHandle()

	// 初始化连接监听器
	tcpAddr, err := net.ResolveTCPAddr("tcp4", Config.GetServerIp())
	Utils.LogErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	Utils.LogErr(err)

	Utils.LogInfo("chat server start listenAndServe！\n")
	i := 0
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			Utils.LogInfo("AcceptTCP error, error = %s", err.Error())
			return
		}
		// conn.SetKeepAlive(true)
		i++
		go Handle.HandleChat(conn, i)
		//		go Handle.Test(i,roomList, conn)
		//		roomChan <- roomList
		Utils.LogInfo("clinet--%d is connected！\n", i)
		time.Sleep(1e8)
	}
	Utils.LogInfo("chat server close！")
}

/**
信号量处理函数
**/
func signalHandle() {
	for {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGUSR1)
		sig := <-ch
		Utils.LogInfo("Signal received(收到信号): %v", sig)
		switch sig {
		case syscall.SIGINT:
			os.Exit(1)
		case syscall.SIGUSR1:
			Utils.LogInfo("重新读取配置文件！")
			Config.ReParseConfig()
		}
	}
}
