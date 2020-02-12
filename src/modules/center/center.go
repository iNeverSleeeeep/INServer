package center

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/lifetime/finalize"
	"INServer/src/modules/cluster"
	"INServer/src/modules/etcmgr"
	"INServer/src/modules/innet"
	"INServer/src/proto/msg"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
)

// Instance 中心服务器的单例
var Instance *Center

const (
	timeout = 2000 * 1E6 // 纳秒
)

type (
	// Center 中心服务器
	Center struct {
		Net           *innet.INNet
		keepAliveTime []int64
	}
)

// New 创建中心服务器
func New() *Center {
	c := new(Center)
	c.Net = innet.New()
	c.keepAliveTime = make([]int64, len(etcmgr.Instance.Servers()))
	c.tickServerState()
	return c
}

// Start 启动中心服务器
func (c *Center) Start() {
	c.Net.Start()
	cluster.SetNode(global.CurrentServerID, &msg.Node{
		NodeState:   msg.NodeState_Running,
		NodeAddress: c.Net.IP,
	})
	c.registerListeners()
	logger.Info(fmt.Sprintf("Server Start Type:%s ID:%d", global.CurrentServerType, global.CurrentServerID))
}

func (c *Center) registerListeners() {
	c.Net.Listen(msg.CMD_NODE_START_NTF, c.HANDLE_NODE_START_NTF)
	c.Net.Listen(msg.CMD_KEEP_ALIVE, c.HANDLE_KEEP_ALIVE)
	c.Net.Listen(msg.CMD_RELOAD_ETC_REQ, c.HANDLE_RELOAD_ETC_REQ)
	c.Net.Listen(msg.CMD_NODES_INFO_NTF, c.HANDLE_NODES_INFO_NTF)
}

// HANDLE_NODE_START_NTF 节点启动
func (c *Center) HANDLE_NODE_START_NTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.NodeStartNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		return
	}
	cluster.SetNode(header.From, &msg.Node{
		NodeState:   msg.NodeState_Unset,
		NodeAddress: ntf.Address,
	})
	c.Net.RefreshNodesAddress()
	c.Net.ResetServer(header.From)
	c.broadcastResetConnectionNTF(header.From)
	c.sendETCSyncNTF(header.From)
}

// HANDLE_NODES_INFO_NTF 每个节点会同步自己的状态过来
func (c *Center) HANDLE_NODES_INFO_NTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.NodesInfoNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		return
	}
	cluster.SetNodes(ntf.Nodes)
	c.Net.RefreshNodesAddress()
	c.pushNodesInfo()
}

func (c *Center) pushNodesInfo() {
	nodes := cluster.GetNodes()
	ntf := &msg.NodesInfoNTF{Nodes: nodes}
	for serverID, node := range nodes {
		if serverID != int(global.CenterID) {
			if node.NodeState == msg.NodeState_Ready || node.NodeState == msg.NodeState_Running {
				c.Net.NotifyServer(msg.CMD_NODES_INFO_NTF, ntf, int32(serverID))
			}
		}
	}
	c.printServerState()
}

// HANDLE_KEEP_ALIVE 心跳消息
func (c *Center) HANDLE_KEEP_ALIVE(header *msg.MessageHeader, buffer []byte) {
	message := &msg.KeepAlive{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		return
	}
	serverID := message.ServerID
	c.keepAliveTime[serverID] = time.Now().UnixNano() + timeout
	info := cluster.GetNode(serverID)
	if info.NodeState == msg.NodeState_Offline {
		info.NodeState = msg.NodeState_Running
		c.pushNodesInfo()
	}
}

// HANDLE_RELOAD_ETC_REQ 重载配置
func (c *Center) HANDLE_RELOAD_ETC_REQ(header *msg.MessageHeader, buffer []byte) {
	defer c.Net.ResponceBytes(header, make([]byte, 1))
	dir, _ := os.Getwd()
	etcmgr.Instance.Load(dir + "/etc")
	c.broadcastETCSyncNTF()
}

func (c *Center) tickServerState() {
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			now := time.Now().UnixNano()
			alloffline := true
			stateDirty := false
			for serverID, info := range cluster.GetNodes() {
				if int32(serverID) == global.CenterID {
					continue
				}
				if info.NodeState == msg.NodeState_Running {
					if c.keepAliveTime[serverID] == 0 {
						c.keepAliveTime[serverID] = now
					}
					if c.keepAliveTime[serverID] < now {
						info.NodeState = msg.NodeState_Offline
						serverType := etcmgr.Instance.GetServerType(int32(serverID))
						logger.Info(fmt.Sprintf("Node Offline ID:%d Type:%s", serverID, serverType))
						stateDirty = true
					}
					if info.NodeState != msg.NodeState_Offline {
						alloffline = false
					}
				}
			}
			if stateDirty {
				c.pushNodesInfo()
				if alloffline {
					finalize.Stop <- true
				}
			}
		}
	}()
}

func (c *Center) broadcastETCSyncNTF() {
	ntf := etcmgr.Instance.GenerateETCSyncNTF()
	for serverID, info := range cluster.GetNodes() {
		if int32(serverID) != global.CenterID {
			if info.NodeState == msg.NodeState_Ready || info.NodeState == msg.NodeState_Running {
				c.Net.NotifyServer(msg.CMD_ETC_SYNC_NTF, ntf, int32(serverID))
			}
		}
	}
}

func (c *Center) sendETCSyncNTF(serverID int32) {
	ntf := etcmgr.Instance.GenerateETCSyncNTF()
	c.Net.NotifyServer(msg.CMD_ETC_SYNC_NTF, ntf, int32(serverID))
}

func (c *Center) broadcastResetConnectionNTF(server int32) {
	ntf := &msg.ResetConnectionNTF{
		ServerID: server,
	}
	for serverID, info := range cluster.GetNodes() {
		if int32(serverID) != global.CenterID {
			if info.NodeState != msg.NodeState_Unset {
				c.Net.NotifyServer(msg.CMD_RESET_CONNECTION_NTF, ntf, int32(serverID))
			}
		}
	}
}

func (c *Center) printServerState() {
	states := "[ServerState] "
	for _, svr := range etcmgr.Instance.Servers() {
		if svr.ServerID == 0 {
			continue
		}
		info := cluster.GetNode(svr.ServerID)
		running := info.NodeState == msg.NodeState_Running
		ready := info.NodeState == msg.NodeState_Ready
		offline := info.NodeState == msg.NodeState_Offline
		if running {
			states = states + fmt.Sprintf("%d%s ", svr.ServerID, "O")
		} else if ready {
			states = states + fmt.Sprintf("%d%s ", svr.ServerID, "R")
		} else if offline {
			states = states + fmt.Sprintf("%d%s ", svr.ServerID, "?")
		} else {
			states = states + fmt.Sprintf("%d%s ", svr.ServerID, "X")
		}
	}
	logger.Info(states)
}
