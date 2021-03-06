package network

import (
	"fmt"
	"log"
	"mynet/base"
	"mynet/rpc"
	"net"
	"sync"
	"sync/atomic"
)

type IServerSocket interface {
	ISocket

	AssignClientId() uint32
	GetClientById(uint32) *ServerSocketClient
	LoadClient() *ServerSocketClient
	AddClinet(*net.TCPConn, string, int) *ServerSocketClient
	DelClinet(*ServerSocketClient) bool
	StopClient(uint32)
}

type ServerSocket struct {
	Socket
	m_nClientCount int
	m_nMaxClients  int
	m_nMinClients  int
	m_nIdSeed      uint32
	m_ClientList   map[uint32]*ServerSocketClient
	m_ClientLocker *sync.RWMutex
	m_Listen       *net.TCPListener
	m_Lock         sync.Mutex
}

type ClientChan struct {
	pClient *ServerSocketClient
	state   int
	id      int
}

type WriteChan struct {
	buff []byte
	id   int
}

func (this *ServerSocket) Init(ip string, port int) bool {
	this.Socket.Init(ip, port)
	this.m_ClientList = make(map[uint32]*ServerSocketClient)
	this.m_ClientLocker = &sync.RWMutex{}
	this.m_sIP = ip
	this.m_nPort = port
	return true
}
func (this *ServerSocket) Start() bool {
	this.m_bShuttingDown = false

	if this.m_sIP == "" {
		this.m_sIP = "127.0.0.1"
	}

	var strRemote = fmt.Sprintf("%s:%d", this.m_sIP, this.m_nPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
	if err != nil {
		log.Fatalf("%v", err)
	}
	ln, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Fatalf("%v", err)
		return false
	}

	fmt.Printf("启动监听，等待链接！\n")

	this.m_Listen = ln
	//延迟，监听关闭
	//defer ln.Close()
	this.m_nState = SSF_ACCEPT
	go this.Run()
	return true
}

//初始化种子
func (this *ServerSocket) AssignClientId() uint32 {
	return atomic.AddUint32(&this.m_nIdSeed, 1)
}

//
func (this *ServerSocket) GetClientById(id uint32) *ServerSocketClient {
	this.m_ClientLocker.RLock()
	client, exist := this.m_ClientList[id]
	this.m_ClientLocker.RUnlock()
	if exist == true {
		return client
	}

	return nil
}

func (this *ServerSocket) AddClinet(tcpConn *net.TCPConn, addr string, connectType int) *ServerSocketClient {
	pClient := this.LoadClient()
	if pClient != nil {
		pClient.Init("", 0)
		pClient.m_pServer = this
		pClient.m_ReceiveBufferSize = this.m_ReceiveBufferSize
		pClient.m_MaxReceiveBufferSize = this.m_MaxReceiveBufferSize
		//客户端id自增
		pClient.m_ClientId = this.AssignClientId()
		pClient.m_sIP = addr
		pClient.SetConnectType(connectType)
		pClient.SetTcpConn(tcpConn)
		this.m_ClientLocker.Lock()
		this.m_ClientList[pClient.m_ClientId] = pClient
		this.m_ClientLocker.Unlock()
		pClient.Start()
		this.m_nClientCount++
		return pClient
	} else {
		log.Printf("%s", "无法创建客户端连接对象")
	}
	return nil
}

func (this *ServerSocket) DelClinet(pClient *ServerSocketClient) bool {
	this.m_ClientLocker.Lock()
	delete(this.m_ClientList, pClient.m_ClientId)
	this.m_ClientLocker.Unlock()
	return true
}

func (this *ServerSocket) StopClient(id uint32) {
	pClinet := this.GetClientById(id)
	if pClinet != nil {
		pClinet.Stop()
	}
}

func (this *ServerSocket) LoadClient() *ServerSocketClient {
	s := &ServerSocketClient{}
	return s
}

//停止该server
func (this *ServerSocket) Stop() bool {
	if this.m_bShuttingDown {
		return true
	}

	this.m_bShuttingDown = true
	this.m_nState = SSF_SHUT_DOWN
	return true
}

func (this *ServerSocket) Send(head rpc.RpcHead, buff []byte) int {
	pClient := this.GetClientById(head.SocketId)
	if pClient != nil {
		pClient.Send(head, buff)
	}
	return 0
}

func (this *ServerSocket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	pClient := this.GetClientById(head.SocketId)
	if pClient != nil {
		pClient.Send(head, base.SetTcpEnd(rpc.Marshal(head, funcName, params...)))
	}
}

func (this *ServerSocket) Restart() bool {
	return true
}

func (this *ServerSocket) Connect() bool {
	return true
}

func (this *ServerSocket) Disconnect(bool) bool {
	return true
}

func (this *ServerSocket) OnNetFail(int) {
}

func (this *ServerSocket) Close() {
	defer this.m_Listen.Close()
	this.Clear()
}

func (this *ServerSocket) Run() bool {
	for {
		tcpConn, err := this.m_Listen.AcceptTCP()
		handleError(err)
		if err != nil {
			return false
		}

		fmt.Printf("客户端：%s已连接！\n", tcpConn.RemoteAddr().String())
		//延迟，关闭链接
		//defer tcpConn.Close()
		this.handleConn(tcpConn, tcpConn.RemoteAddr().String())
	}
}

func (this *ServerSocket) handleConn(tcpConn *net.TCPConn, addr string) bool {
	if tcpConn == nil {
		return false
	}

	pClient := this.AddClinet(tcpConn, addr, this.m_nConnectType)
	if pClient == nil {
		return false
	}

	return true
}
