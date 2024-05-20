package ziface


type IPackage interface {
	// 封包操作
	Pack(msg IMessage) ([]byte, error)
	// 解包操作
	Unpack([]byte) (IMessage, error)
	// 获取包头部长度
	GetHeadLen() uint32
}

