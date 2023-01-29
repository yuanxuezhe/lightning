// 格式化TCP包数据
// 数据格式
// --------------------
// | lenMsglen | data |
// --------------------
// lenMsglen 长度消息的字节长度
// data 消息体
package network

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

var (
	lenMsglen int  = 4
	bigEndian bool = true
)

// 设置消息体最大长度，最小长度
func GetMsgLen() (minMsgLen, maxMsgLen uint32, err error) {
	var max uint32 // uint32长度消息体最长4字节
	if lenMsglen > 4 {
		err = errors.New("消息长度记录字节最多不能超过4！")
	}

	max = 1<<(8*lenMsglen) - 1

	minMsgLen = 1
	maxMsgLen = max

	return
}

// 格式化TCP数据包 发送消息
func SendMsg(conn net.Conn, msg []byte) error {
	// 取消息体长度范围
	minMsgLen, maxMsgLen, err := GetMsgLen()
	if err != nil {
		return err
	}
	msgLen := uint32(len(msg))

	// 校验消息体范围
	if msgLen > maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < minMsgLen {
		return errors.New("message too short")
	}

	ansMsg := make([]byte, uint32(lenMsglen)+msgLen)

	// 根据大端序或小端序，将消息长度写入TCP消息串
	if bigEndian {
		binary.BigEndian.PutUint32(ansMsg, msgLen)
	} else {
		binary.LittleEndian.PutUint32(ansMsg, msgLen)
	}

	// 写入消息实体数据
	copy(ansMsg[lenMsglen:], msg)

	_, err = conn.Write(ansMsg)
	if err != nil {
		return err
	}
	return nil
}

// 解析获取消息实体
func ReadMsg(conn net.Conn) ([]byte, error) {
	// 取数据流消息体长度
	bufMsgLen := make([]byte, lenMsglen)
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}

	// 解析长度
	var msgLen uint32
	if bigEndian {
		msgLen = binary.BigEndian.Uint32(bufMsgLen)
	} else {
		msgLen = binary.LittleEndian.Uint32(bufMsgLen)
	}

	// 取消息体长度范围
	minMsgLen, maxMsgLen, err := GetMsgLen()
	if err != nil {
		return nil, err
	}

	// 校验消息体范围
	if msgLen > maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < minMsgLen {
		return nil, errors.New("message too short")
	}

	// 获取数据
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, err
	}

	return msgData, nil
}

// 解析获取消息实体
func ReadMsgs(conn net.Conn) ([]string, error) {
	// 取数据流消息体长度
	list := []string{}

	for {
		buff, err := ReadMsg(conn)
		if err != nil {
			if err.Error() == "EOF" && len(list) > 0 {
				return list, nil
			} else {
				return nil, err
			}

		} else {
			list = append(list, string(buff))
		}
	}
}
