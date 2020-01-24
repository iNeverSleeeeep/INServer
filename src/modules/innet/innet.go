package innet

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/profiler"
	"INServer/src/common/protect"
	"INServer/src/proto/msg"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

var sequence uint64 = 0

var udpconn *net.UDPConn

type (
	responce struct {
		c       chan []byte
		timeout int64
	}

	INNet struct {
		responces map[uint64]*responce
		listeners map[msg.Command]chan *msg.Message
		sender    *sender
		address   *address
		retry     *retry
		receiver  *receiver
		conn      *net.UDPConn
		IP        []byte

		gates    []int32
		database int32
	}
)

func New() *INNet {
	if udpconn != nil {
		logger.Fatal("innet new wrice not allow")
	}
	innet := new(INNet)
	innet.IP = getIP()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: int(global.ServerID) + sendport})
	if err != nil {
		logger.Fatal(err)
	}
	udpconn = conn
	innet.conn = conn
	innet.database = global.InvalidServerID
	innet.responces = make(map[uint64]*responce)
	innet.listeners = make(map[msg.Command]chan *msg.Message)
	innet.address = newAddress(innet)
	innet.retry = newRetry(innet)
	innet.receiver = newReceiver(innet)
	innet.sender = newSender(innet)
	innet.startTickRequestTimeout()
	return innet
}

func (n *INNet) Start() {
	n.receiver.start()
}

func (n *INNet) Request(command msg.Command, req proto.Message, resp proto.Message) error {
	sequence++
	err := n.sendMessage(command, sequence, req, global.InvalidServerID)
	if err != nil {
		return err
	}
	buffer := make(chan []byte)
	n.responces[sequence] = &responce{c: buffer, timeout: time.Now().Unix() + 10}
	buf := <-buffer
	delete(n.responces, sequence)
	if buf != nil {
		return proto.Unmarshal(buf, resp)
	} else {
		return nil
	}
}

// RequestClientBytes 发送源头为客户端的消息
func (n *INNet) RequestClientBytes(command msg.Command, uuid string, bytes []byte) ([]byte, error) {
	sequence++
	err := n.sendClientBytes(command, sequence, uuid, bytes, global.InvalidServerID)
	if err != nil {
		return nil, err
	}
	buffer := make(chan []byte)
	n.responces[sequence] = &responce{c: buffer, timeout: time.Now().Unix() + 10}
	buf := <-buffer
	delete(n.responces, sequence)
	return buf, nil
}

// RequestBytes 发送
func (n *INNet) RequestBytes(command msg.Command, bytes []byte) ([]byte, error) {
	sequence++
	err := n.sendBytes(command, sequence, bytes, global.InvalidServerID)
	if err != nil {
		return nil, err
	}
	buffer := make(chan []byte)
	n.responces[sequence] = &responce{c: buffer, timeout: time.Now().Unix() + 10}
	buf := <-buffer
	delete(n.responces, sequence)
	return buf, nil
}

func (n *INNet) Responce(header *msg.MessageHeader, message proto.Message) error {
	return n.sendMessage(msg.Command_RESP, header.Sequence, message, header.From)
}

func (n *INNet) ResponceBytes(header *msg.MessageHeader, bytes []byte) error {
	return n.sendBytes(msg.Command_RESP, header.Sequence, bytes, header.From)
}

func (n *INNet) Notify(command msg.Command, message proto.Message) error {
	sequence++
	return n.sendMessage(command, sequence, message, global.InvalidServerID)
}

// NotifyClientBytes 发送源头为客户端的消息
func (n *INNet) NotifyClientBytes(command msg.Command, uuid string, bytes []byte) error {
	sequence++
	return n.sendClientBytes(command, sequence, uuid, bytes, global.InvalidServerID)
}

func (n *INNet) NotifyBytes(command msg.Command, bytes []byte) error {
	sequence++
	return n.sendBytes(command, sequence, bytes, global.InvalidServerID)
}

func (n *INNet) NotifyServer(command msg.Command, message proto.Message, serverID int32) error {
	sequence++
	return n.sendMessage(command, sequence, message, serverID)
}

// TODO 缺少取消监听
func (n *INNet) Listen(command msg.Command, listener func(header *msg.MessageHeader, buffer []byte)) {
	pkgChan := make(chan *msg.Message)
	n.listeners[command] = pkgChan
	go func() {
		for {
			select {
			case pkg := <-pkgChan:
				go func() {
					defer protect.CatchPanic()
					listener(pkg.Header, pkg.Buffer)
				}()
			}
		}
	}()
}

func (n *INNet) Gates() []int32 {
	return n.gates
}

// GetGateAddress 网关内网地址
func (n *INNet) GetGateAddress(serverID int32) ([]byte, int) {
	if info, ok := n.address.servers[serverID]; ok {
		return info.info.Address, int(global.Servers[serverID].ServerConfig.GateConfig.Port)
	}
	return nil, 0
}

// GetGatePublicAddress 网关公网地址
func (n *INNet) GetGatePublicAddress(serverID int32) (string, int) {
	return global.Servers[serverID].ServerConfig.GateConfig.Address, int(global.Servers[serverID].ServerConfig.GateConfig.Port)
}

func (n *INNet) AddServers(servers []*msg.ServerInfo) {
	n.address.addServerList(servers)
	n.refreshRunningServers()
}

func (n *INNet) ResetServers(servers []*msg.ServerInfo) {
	for _, svr := range servers {
		n.address.resetServer(svr.ServerID)
		n.receiver.resetServer(svr.ServerID)
		n.retry.resetServer(svr.ServerID)
	}
	n.refreshRunningServers()
}

func (n *INNet) ResetServer(svr *msg.ServerInfo) {
	n.address.resetServer(svr.ServerID)
	n.receiver.resetServer(svr.ServerID)
	n.retry.resetServer(svr.ServerID)
	n.refreshRunningServers()
}

func (n *INNet) IsServerRunning(serverID int32) bool {
	for _, info := range n.address.servers {
		if info.info.ServerID == serverID {
			return info.info.State == msg.ServerState_Running
		}
	}
	return false
}

func (n *INNet) sendMessage(command msg.Command, sequence uint64, message proto.Message, serverID int32) error {
	var err error
	if buffer, err := proto.Marshal(message); err == nil {
		return n.sendBytes(command, sequence, buffer, serverID)
	}
	return err
}

func (n *INNet) sendBytes(command msg.Command, sequence uint64, bytes []byte, serverID int32) error {
	header := &msg.MessageHeader{
		Command:  command,
		Sequence: sequence,
		From:     global.ServerID,
	}
	return n.send(header, bytes, serverID)
}

func (n *INNet) sendClientBytes(command msg.Command, sequence uint64, uuid string, bytes []byte, serverID int32) error {
	header := &msg.MessageHeader{
		Command:    command,
		Sequence:   sequence,
		From:       global.ServerID,
		PlayerUUID: uuid,
	}
	return n.send(header, bytes, serverID)
}

func (n *INNet) send(header *msg.MessageHeader, bytes []byte, serverID int32) error {
	buf, err := proto.Marshal(&msg.Message{Header: header, Buffer: bytes})
	if err == nil {
		var svr *server
		if serverID != global.InvalidServerID {
			svr = n.address.getByServerID(serverID)
		} else {
			svr = n.address.getByCommand(header.Command)
			if svr == nil {
				return errors.New("Server Not Found By Command:" + header.Command.String())
			}
		}
		if header.Command == msg.Command_SERVER_STATE {
			err = n.sender.sendstate(svr, buf)
		} else {
			err = n.sender.send(svr, buf)
		}
	}
	return err
}

func (n *INNet) handleMessage(message *msg.Message) {
	if listener, ok := n.listeners[message.Header.Command]; ok {
		sequence := message.Header.Sequence
		profiler.BeginSampleMessage(sequence)
		listener <- message
		profiler.EndSampleMessage(sequence)
	} else if resp, ok := n.responces[message.Header.Sequence]; ok {
		resp.c <- message.Buffer
	} else {
		logger.Debug("handleMessage error:" + message.Header.String())
	}
}

func (n *INNet) refreshRunningServers() {
	n.gates = make([]int32, 0)
	n.database = global.InvalidServerID
	for _, info := range n.address.servers {
		if info.info.State == msg.ServerState_Running {
			// FIXME 这里有个异常 runtime error: index out of range
			if len(global.Servers) < int(info.info.ServerID) {
				logger.Error(fmt.Sprintf("len:%dm serverID:%d", len(global.Servers), info.info.ServerID))
			}
			serverType := global.Servers[int(info.info.ServerID)].ServerType
			if serverType == global.GateServer {
				n.gates = append(n.gates, info.info.ServerID)
			} else if serverType == global.DatabaseServer {
				n.database = info.info.ServerID
			}
		}
	}
}

func (n *INNet) startTickRequestTimeout() {
	go func() {
		for {
			chans := make([]chan []byte, 0)
			now := time.Now().Unix()
			for _, resp := range n.responces {
				if resp.timeout <= now {
					chans = append(chans, resp.c)
				}
			}
			for _, c := range chans {
				c <- nil
			}
			time.Sleep(time.Second)
		}
	}()
}
