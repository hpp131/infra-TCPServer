package znet

/*
该模块意在解决TCP粘包问题，通过在应用程序中规范Message实体，并且实现拆包和解包逻辑来解决粘包问题
核心功能是实现封包和解包(Pack, Unpack)
*/

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"tcpserver/util"
	"tcpserver/ziface"
)

func NewDataPackage() *DataPackage {
	return &DataPackage{
	}
}

type DataPackage struct {
	Head []byte
	Data []byte
	
}

func (dp *DataPackage) GetHeadLen() uint32 {
	// Message.len int32 = 4bytes
	// Message.MsgID int32 = 4bytes
	return 8
}

func (dp *DataPackage) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen());err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgID());err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData());err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

func (dp *DataPackage) Unpack(data []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Len);err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.MsgID);err != nil {
		return nil, err
	}
	if util.Globalobject.MaxPackageSize > 0 && msg.Len > util.Globalobject.MaxPackageSize {
		fmt.Printf("MaxPackageSize %d, msg.Len %d", util.Globalobject.MaxPackageSize, msg.Len)
		return nil, errors.New("message is too large")
	}
	return msg, nil
}
