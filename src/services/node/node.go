package node

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/proto/msg"
	"INServer/src/services/cluster"
	"INServer/src/services/etcmgr"
	"INServer/src/services/innet"
	"time"

	"github.com/golang/protobuf/proto"
)

var (
	Instance *Node
	Net      *innet.INNet
)

type (
	Node struct {
	}
)

func New() *Node {
	n := new(Node)
	Net = innet.New()
	n.registerListeners()
	Net.Start()
	return n
}

// Prepare 节点进入Ready状态 这个状态之后可以收到集群的状态信息了
func (n *Node) Prepare() {
	ntf := &msg.NodesInfoNTF{}
	ntf.Nodes = make([]*msg.Node, 0)
	for index := 0; index < len(etcmgr.Instance.Servers()); index++ {
		ntf.Nodes = append(ntf.Nodes, &msg.Node{
			NodeState:   msg.NodeState_Unset,
			NodeAddress: nil,
		})
	}
	ntf.Nodes[global.CurrentServerID] = &msg.Node{
		NodeState:   msg.NodeState_Ready,
		NodeAddress: Net.IP,
	}
	Net.NotifyServer(msg.CMD_NODES_INFO_NTF, ntf, global.CenterID)
}

// Start 节点进入Running状态 工作状态
func (n *Node) Start() {
	ntf := &msg.NodesInfoNTF{}
	ntf.Nodes = make([]*msg.Node, 0)
	for index := 0; index < len(etcmgr.Instance.Servers()); index++ {
		ntf.Nodes = append(ntf.Nodes, &msg.Node{
			NodeState:   msg.NodeState_Unset,
			NodeAddress: nil,
		})
	}
	ntf.Nodes[global.CurrentServerID] = &msg.Node{
		NodeState:   msg.NodeState_Running,
		NodeAddress: Net.IP,
	}
	Net.NotifyServer(msg.CMD_NODES_INFO_NTF, ntf, global.CenterID)
	n.keepAlive()
}

func (n *Node) registerListeners() {
	Net.Listen(msg.CMD_NODES_INFO_NTF, n.HANDLE_NODES_INFO_NTF)
	Net.Listen(msg.CMD_ETC_SYNC_NTF, etcmgr.Instance.HANDLE_ETC_SYNC_NTF)
	Net.Listen(msg.CMD_RESET_CONNECTION_NTF, n.HANDLE_RESET_CONNECTION_NTF)
}

func (n *Node) HANDLE_NODES_INFO_NTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.NodesInfoNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		logger.Error(err)
		return
	}
	cluster.SetNodes(ntf.Nodes)
	Net.RefreshNodesAddress()
}

func (n *Node) HANDLE_RESET_CONNECTION_NTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.ResetConnectionNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		logger.Error(err)
		return
	}
	Net.ResetServer(ntf.ServerID)
}

func (n *Node) keepAlive() {
	go func() {
		for {
			info := &msg.KeepAlive{
				ServerID: global.CurrentServerID,
			}
			Net.NotifyServer(msg.CMD_KEEP_ALIVE, info, 0)
			time.Sleep(time.Second)
		}
	}()
}
