package netgate

import (
	"github.com/golang/protobuf/proto"
	"mynet/base"
	"mynet/rpc"
	"mynet/server/message"
	"strings"
)

var (
	A_C_RegisterResponse = strings.ToLower("A_C_RegisterResponse")
	A_C_LoginResponse    = strings.ToLower("A_C_LoginResponse")
)

func SendToClient(socketId uint32, packet proto.Message) {
	SERVER.GetServer().Send(rpc.RpcHead{SocketId: socketId}, base.SetTcpEnd(message.Encode(packet)))
}

//此函数为消息队列订阅的回调函数，当account，world像消息队列里放消息时就会触发此函数
func DispatchPacket(id uint32, buff []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	rpcPacket, head := rpc.Unmarshal(buff)
	switch head.DestServerType {
	//如果是发到网关的
	case rpc.SERVICE_GATESERVER:
		bitstream := base.NewBitStream(rpcPacket.RpcBody, len(rpcPacket.RpcBody))
		buff := message.EncodeEx(rpcPacket.FuncName, rpc.UnmarshalPB(bitstream))
		//account到client的注册登录的回应转发给client
		if rpcPacket.FuncName == A_C_RegisterResponse || rpcPacket.FuncName == A_C_LoginResponse {
			SERVER.GetServer().Send(rpc.RpcHead{SocketId: head.SocketId}, base.SetTcpEnd(buff))
		} else {
			socketId := SERVER.GetPlayerMgr().GetSocket(head.Id)
			SERVER.GetServer().Send(rpc.RpcHead{SocketId: socketId}, base.SetTcpEnd(buff))
		}
	default:
		//转发
		SERVER.GetCluster().Send(head, buff)
	}

	return true
}
