package cluster

import (
	"INServer/src/proto/msg"
	"INServer/src/services/etcmgr"
)

var (
	nodes []*msg.Node
)

// SetNodes 设置集群信息
func SetNodes(nodesArray []*msg.Node) {
	for serverID, node := range nodesArray {
		SetNode(int32(serverID), node)
	}
	RefreshRunning()
	RefreshRunningZones()
}

// SetNode 设置单个节点信息
func SetNode(serverID int32, node *msg.Node) {
	for len(nodes) < int(serverID+1) {
		nodes = append(nodes, &msg.Node{})
	}
	nodes[serverID].NodeAddress = node.NodeAddress
	if node.NodeState == msg.NodeState_Unset {
		return
	} else if node.NodeState == msg.NodeState_Ready && nodes[serverID].NodeState == msg.NodeState_Running {
		// 因为消息乱序处理的问题，有可能先处理runining后处理ready 这个时候状态不要切换错了
		return
	}
}

// GetNode 取得单个节点信息
func GetNode(serverID int32) *msg.Node {
	if len(nodes) < int(serverID+1) {
		return &msg.Node{}
	}
	if nodes[serverID] != nil {
		return nodes[serverID]
	}
	return &msg.Node{}
}

// GetNodeState 取得节点状态
func GetNodeState(serverID int32) msg.NodeState {
	if len(nodes) < int(serverID+1) {
		return msg.NodeState_Unset
	}
	if nodes[serverID] == nil {
		return msg.NodeState_Unset
	}
	return nodes[serverID].NodeState
}

// GetNodes 取得集群信息
func GetNodes() []*msg.Node {
	return nodes
}

// GetGatePublicAddress 网关公网地址
func GetGatePublicAddress(serverID int32) (string, int, int) {
	servers := etcmgr.Instance.Servers()
	ip := servers[int(serverID)].ServerConfig.GateConfig.Address
	port := int(servers[int(serverID)].ServerConfig.GateConfig.Port)
	webport := int(servers[int(serverID)].ServerConfig.GateConfig.WebPort)
	return ip, port, webport
}

func init() {
	nodes = make([]*msg.Node, 0)
}
