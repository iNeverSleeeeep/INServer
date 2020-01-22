package node

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/modules/innet"
	"INServer/src/proto/msg"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
)

var (
	Instance *Node
)

type (
	Node struct {
		Net *innet.INNet
	}
)

func New() *Node {
	n := new(Node)
	n.Net = innet.New()
	n.registerListeners()
	n.Net.Start()
	return n
}

func (n *Node) Prepare() {
	info := &msg.ServerInfo{
		ServerID: global.ServerID,
		Address:  n.Net.IP,
		State:    msg.ServerState_Start,
	}
	req := &msg.ServerStateReq{
		Info: info,
	}
	resp := &msg.ServerStateResp{}
	err := n.Net.Request(msg.Command_SERVER_STATE, req, resp)
	if err != nil {
		logger.Fatal(err)
	}

	global.ServerType = resp.ServerType
	global.ServerConfig = resp.ServerConfig
	global.Servers = resp.Servers
	global.Zones = resp.Zones

	if global.ServerType == global.InvalidServer {
		logger.Debug(resp.Message)
		os.Exit(0)
	}

	logger.Info(fmt.Sprintf("Server Start Type:%s ID:%d", resp.ServerType, global.ServerID))
}

func (n *Node) Start() {
	info := &msg.ServerInfo{
		ServerID: global.ServerID,
		Address:  n.Net.IP,
		State:    msg.ServerState_Running,
	}
	req := &msg.ServerStateReq{
		Info: info,
	}
	resp := &msg.ServerStateResp{}
	err := n.Net.Request(msg.Command_SERVER_STATE, req, resp)
	if err != nil {
		logger.Debug(err)
	}
	n.keepAlive()
}

func (n *Node) registerListeners() {
	n.Net.Listen(msg.Command_SERVER_LIST_CHANGED, n.onServerListChanged)
	n.Net.Listen(msg.Command_SERVER_RESET, n.onServerReset)
	n.Net.Listen(msg.Command_UPDATE_ETC_NTF, n.onUpdateETC)
}

func (n *Node) onServerListChanged(header *msg.MessageHeader, buffer []byte) {
	message := &msg.ServerInfoList{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		return
	}
	n.Net.AddServers(message.Servers)
}

func (n *Node) onServerReset(header *msg.MessageHeader, buffer []byte) {
	message := &msg.ServerInfoList{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		return
	}
	n.Net.AddServers(message.Servers)
	n.Net.ResetServers(message.Servers)
}

func (n *Node) onUpdateETC(header *msg.MessageHeader, buffer []byte) {
	message := &msg.UpdateETC{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		return
	}
	if message.ServerConfig != nil {
		global.ServerConfig = message.ServerConfig
	}
	if message.Servers != nil {
		global.Servers = message.Servers
	}
	if message.Zones != nil {
		global.Zones = message.Zones
	}
}

func (n *Node) keepAlive() {
	go func() {
		for {
			info := &msg.KeepAlive{
				ServerID: global.ServerID,
			}
			n.Net.NotifyServer(msg.Command_KEEP_ALIVE, info, 0)
			time.Sleep(time.Second)
		}
	}()
}
