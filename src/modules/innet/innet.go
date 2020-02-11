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
		listeners map[msg.CMD]chan *msg.Message
		sender    *sender
		address   *address
		retry     *retry
		receiver  *receiver
		conn      *net.UDPConn
		IP        []byte
	}
)

func New() *INNet {
	if udpconn != nil {
		logger.Fatal("innet new wrice not allow")
	}
	innet := new(INNet)
	innet.IP = getIP()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: int(global.CurrentServerID) + sendport})
	if err != nil {
		logger.Fatal(err)
	}
	udpconn = conn
	innet.conn = conn
	innet.responces = make(map[uint64]*responce)
	innet.listeners = make(map[msg.CMD]chan *msg.Message)
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

func (n *INNet) Request(command msg.CMD, req proto.Message, resp proto.Message) error {
	return n.RequestServer(command, req, resp, global.InvalidServerID)
}

func (n *INNet) RequestServer(command msg.CMD, req proto.Message, resp proto.Message, serverID int32) error {
	sequence++
	err := n.sendMessage(command, sequence, req, serverID)
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
func (n *INNet) RequestClientBytes(command msg.CMD, uuid string, bytes []byte) ([]byte, error) {
	return n.RequestClientBytesToServer(command, uuid, bytes, global.InvalidServerID)
}

// RequestClientBytesToServer 发送源头为客户端的消息
func (n *INNet) RequestClientBytesToServer(command msg.CMD, uuid string, bytes []byte, serverID int32) ([]byte, error) {
	sequence++
	err := n.sendClientBytes(command, sequence, uuid, bytes, serverID)
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
func (n *INNet) RequestBytes(command msg.CMD, bytes []byte) ([]byte, error) {
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
	return n.sendMessage(msg.CMD_RESP, header.Sequence, message, header.From)
}

func (n *INNet) ResponceBytes(header *msg.MessageHeader, bytes []byte) error {
	return n.sendBytes(msg.CMD_RESP, header.Sequence, bytes, header.From)
}

func (n *INNet) Notify(command msg.CMD, message proto.Message) error {
	sequence++
	return n.sendMessage(command, sequence, message, global.InvalidServerID)
}

// NotifyClientBytes 发送源头为客户端的消息
func (n *INNet) NotifyClientBytes(command msg.CMD, uuid string, bytes []byte) error {
	sequence++
	return n.sendClientBytes(command, sequence, uuid, bytes, global.InvalidServerID)
}

// NotifyClientBytesToServer 发送源头为客户端的消息
func (n *INNet) NotifyClientBytesToServer(command msg.CMD, uuid string, bytes []byte, serverID int32) error {
	sequence++
	return n.sendClientBytes(command, sequence, uuid, bytes, serverID)
}

func (n *INNet) NotifyBytes(command msg.CMD, bytes []byte) error {
	sequence++
	return n.sendBytes(command, sequence, bytes, global.InvalidServerID)
}

// SendNodeStartNTF 服务器启动通知 只在服务器启动时发送
func (n *INNet) SendNodeStartNTF() {
	sequence++
	ntf := &msg.NodeStartNTF{Address: n.IP}
	if buffer, err := proto.Marshal(ntf); err == nil {
		header := &msg.MessageHeader{
			Command:  msg.CMD_NODE_START_NTF,
			Sequence: sequence,
			From:     global.CurrentServerID,
		}
		svr := n.address.getByServerID(global.CenterID)
		bytes, _ := proto.Marshal(&msg.Message{Header: header, Buffer: buffer})
		n.sender.start(svr, bytes)
	}
}

func (n *INNet) NotifyServer(command msg.CMD, message proto.Message, serverID int32) error {
	sequence++
	return n.sendMessage(command, sequence, message, serverID)
}

// TODO 缺少取消监听
func (n *INNet) Listen(command msg.CMD, listener func(header *msg.MessageHeader, buffer []byte)) {
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

// RefreshNodesAddress 刷新地址
func (n *INNet) RefreshNodesAddress() {
	n.address.refresh()
}

// ResetServer 重置服务器连接状态
func (n *INNet) ResetServer(serverID int32) {
	n.address.resetServer(serverID)
	n.receiver.resetServer(serverID)
	n.retry.resetServer(serverID)
}

func (n *INNet) sendMessage(command msg.CMD, sequence uint64, message proto.Message, serverID int32) error {
	var err error
	if buffer, err := proto.Marshal(message); err == nil {
		return n.sendBytes(command, sequence, buffer, serverID)
	}
	return err
}

func (n *INNet) sendBytes(command msg.CMD, sequence uint64, bytes []byte, serverID int32) error {
	header := &msg.MessageHeader{
		Command:  command,
		Sequence: sequence,
		From:     global.CurrentServerID,
	}
	return n.send(header, bytes, serverID)
}

func (n *INNet) sendClientBytes(command msg.CMD, sequence uint64, uuid string, bytes []byte, serverID int32) error {
	header := &msg.MessageHeader{
		Command:    command,
		Sequence:   sequence,
		From:       global.CurrentServerID,
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
			if svr == nil {
				errstr := fmt.Sprintf("Server Not Found By ServerID:%d", serverID)
				logger.Error(errstr)
				return errors.New(errstr)
			}
		} else {
			svr = n.address.getByCommand(header.Command)
			if svr == nil {
				logger.Error("Server Not Found By Command:" + header.Command.String())
				return errors.New("Server Not Found By Command:" + header.Command.String())
			}
		}
		err = n.sender.send(svr, buf)
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
