package ziface

type IConnManager interface {
	AddConn(conn IConnection)
	DeleteConn(conn IConnection)
	GetConn(connID uint32) IConnection
	// 获取所有链接的数量
	Len() int
	// 清理所有链接
	ClearAll()
}
