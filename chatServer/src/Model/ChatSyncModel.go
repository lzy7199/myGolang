package Model

import (
	"net"
	"sync"
)

type ChatSyncModel struct {
	M         sync.Mutex
	SocketMap map[*net.TCPConn](*ChatAvatar)    /**保存socket连接的map**/
	NameMap   map[int]map[string](*net.TCPConn) /**保存conn和name对应关系的map，第一层map的key是serverIndex**/
}
