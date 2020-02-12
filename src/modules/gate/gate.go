package gate

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/protect"
	"INServer/src/modules/innet"
	"INServer/src/modules/node"
	"INServer/src/proto/data"
	"INServer/src/proto/msg"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var Instance *Gate
var upgrader = websocket.Upgrader{}

type (
	Gate struct {
		listener *net.TCPListener
		roles    map[string]*session
		serv     *services
	}
)

func New() *Gate {
	g := new(Gate)
	g.roles = make(map[string]*session)
	g.serv = newServices()
	g.serv.start()
	return g
}

func (g *Gate) Start() {
	// TCP
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: int(global.CurrentServerConfig.GateConfig.Port)})
	if err != nil {
		logger.Fatal(err)
	}
	g.listener = listener

	logger.Info("门服务器 启动 监听端口:" + strconv.Itoa(int(global.CurrentServerConfig.GateConfig.Port)))
	go func() {
		for {
			conn, err := g.listener.AcceptTCP()
			if err != nil {
				logger.Debug("listener.Accept 连接错误")
			} else {
				go g.handleConnect(conn)
			}
		}
	}()

	// WebSocket
	if global.CurrentServerConfig.GateConfig.WebPort > 0 {
		http.HandleFunc("/", g.handleWebConnect)
		go http.ListenAndServe(fmt.Sprintf(":%d", global.CurrentServerConfig.GateConfig.WebPort), nil)
		logger.Info("门服务器 监听端口:" + strconv.Itoa(int(global.CurrentServerConfig.GateConfig.WebPort)))
	}

}

func (g *Gate) initMessageHandler() {
	node.Instance.Net.Listen(msg.CMD_UPDATE_ROLE_ADDRESS_NTF, g.HANDLE_UPDATE_ROLE_ADDRESS_NTF)
}

func (g *Gate) HANDLE_UPDATE_ROLE_ADDRESS_NTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.UpdateRoleAddressNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		logger.Error(err)
		return
	}

	if session, ok := g.roles[ntf.RoleUUID]; ok {
		if ntf.Address.Entity != global.InvalidServerID {
			session.info.Address.Entity = ntf.Address.Entity
		}
	}
}

func (g *Gate) handleWebConnect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("upgrade:", err)
		return
	}
	var uuid *string = nil
	defer g.closeConnection(uuid, nil, c)
	defer c.Close()
	defer protect.CatchPanic()
	var buf = make([]byte, 65536)
	current := 0
	for {
		for current < 2 {
			_, msg, err := c.ReadMessage()
			if err != nil {
				logger.Info("read:", err)
				return
			}
			copy(buf[current:], msg[:])
			current = current + len(msg)
		}
		// 等待读取数据
		size := binary.BigEndian.Uint16(buf[:2])
		for (current - 2) < int(size) {
			_, msg, err := c.ReadMessage()
			if err != nil {
				logger.Info("read:", err)
				return
			}
			copy(buf[current:], msg[:])
			current = current + len(msg)
		}

		message := &msg.ClientToGate{}
		err := proto.Unmarshal(buf[2:size+2], message)
		if err != nil {
			logger.Debug("消息解析失败:", c.RemoteAddr())
			continue
		}
		if message.Command == msg.CMD_CONNECT_GATE_REQ {
			connectReq := &msg.ConnectGateReq{}
			err := proto.Unmarshal(message.Request, connectReq)
			if err != nil {
				logger.Info("消息解析失败")
				continue
			}
			if uuid == nil {
				uuid = &connectReq.SessionCert.UUID
			}
			g.handleConnectMessage(uuid, nil, c, message.Sequence, connectReq)
		} else if uuid != nil {
			if role, ok := g.roles[*uuid]; ok {
				if role.info.State == data.SessionState_Online {
					g.handleMessage(role, message)
				} else {
					logger.Info("客户端状态错误:" + c.RemoteAddr().String() + " 当前状态:" + role.info.State.String())
					return
				}
			}
		} else {
			logger.Info("客户端需要先发送Connect协议:" + c.RemoteAddr().String())
			return
		}

		copy(buf[0:], buf[size+2:current])
		current = current - int(size) - 2
	}
}

func (g *Gate) handleConnect(conn *net.TCPConn) {
	var uuid *string = nil
	defer g.closeConnection(uuid, conn, nil)
	defer protect.CatchPanic()
	var buf = make([]byte, 65536)
	current := 0
	for {
		// 等待读取数据长度
		for current < 2 {
			n, err := conn.Read(buf[current:])
			if err != nil {
				logger.Info(err)
				return
			}
			current = current + n
		}

		// 等待读取数据
		size := binary.BigEndian.Uint16(buf[:2])
		for (current - 2) < int(size) {
			n, err := conn.Read(buf[current:])
			if err != nil {
				logger.Info(err)
				return
			}
			current = current + n
		}

		message := &msg.ClientToGate{}
		err := proto.Unmarshal(buf[2:size+2], message)
		if err != nil {
			logger.Debug("消息解析失败:" + conn.RemoteAddr().String())
			continue
		}
		if message.Command == msg.CMD_CONNECT_GATE_REQ {
			connectReq := &msg.ConnectGateReq{}
			err := proto.Unmarshal(message.Request, connectReq)
			if err != nil {
				logger.Info("消息解析失败")
				continue
			}
			if uuid == nil {
				uuid = &connectReq.SessionCert.UUID
			}
			g.handleConnectMessage(uuid, conn, nil, message.Sequence, connectReq)
		} else if uuid != nil {
			if role, ok := g.roles[*uuid]; ok {
				if role.info.State == data.SessionState_Online {
					g.handleMessage(role, message)
				} else {
					logger.Debug("客户端状态错误:" + conn.RemoteAddr().String() + " 当前状态:" + role.info.State.String())
					return
				}
			}
		} else {
			logger.Debug("客户端需要先发送Connect协议:" + conn.RemoteAddr().String())
			return
		}

		copy(buf[0:], buf[size+2:current])
		current = current - int(size) - 2
	}
}

func (g *Gate) closeConnection(uuid *string, conn *net.TCPConn, webconn *websocket.Conn) {
	if conn != nil {
		conn.Close()
	}
	if webconn != nil {
		webconn.Close()
	}
	if uuid != nil {
		if role, ok := g.roles[*uuid]; ok {
			role.conn = nil
			role.webconn = nil
			role.info.State = data.SessionState_OutOfContact
			role.cert.OutOfDateTime = generateCertOutOfDateTime()
		}
	}
}

func (g *Gate) tickOutOfContact() {
	go func() {
		for {
			time.Sleep(time.Second)
			now := time.Now().Unix()
			var outOfDateRoles []string
			for uuid, role := range g.roles {
				if role.info.State == data.SessionState_OutOfContact && role.cert.OutOfDateTime < now {
					outOfDateRoles = append(outOfDateRoles, uuid)
				}
			}
			for _, uuid := range outOfDateRoles {
				delete(g.roles, uuid)
			}
		}
	}()
}

func (g *Gate) handleConnectMessage(uuid *string, conn *net.TCPConn, webconn *websocket.Conn, sequence uint64, message *msg.ConnectGateReq) {
	connectResp := &msg.ConnectGateResp{}
	resp := &msg.GateToClient{}
	resp.Command = msg.CMD_RESP
	resp.Sequence = sequence
	if conn != nil {
		defer g.sendResp(conn, resp)
	} else {
		defer g.sendWebResp(webconn, resp)
	}
	role, ok := g.roles[*uuid]
	if ok == false {
		addr := ""
		if conn != nil {
			addr = conn.RemoteAddr().String()
		} else if webconn != nil {
			addr = webconn.RemoteAddr().String()
		}
		logger.Debug("拒绝连接，没有数据:" + addr + " uuid:" + *uuid)
		return
	} else if message.SessionCert.Key != role.cert.Key {
		addr := ""
		if conn != nil {
			addr = conn.RemoteAddr().String()
		} else if webconn != nil {
			addr = webconn.RemoteAddr().String()
		}
		logger.Info(fmt.Sprintf("客户端秘钥错误 addr:%d uuid:%s client:%s server:%s", addr, *uuid, message.SessionCert.Key, role.cert.Key))
		return
	} else {
		role.conn = conn
		role.webconn = webconn
		role.generateNewCertKey()
		role.info.State = data.SessionState_Online
		role.info.RoleUUID = *uuid
		role.info.Address = &data.RoleAddress{
			Gate:   global.CurrentServerID,
			Entity: global.InvalidServerID,
		}
		ntf := &msg.UpdateRoleAddressNTF{
			RoleUUID: *uuid,
			Address:  role.info.Address,
		}
		node.Instance.Net.Notify(msg.CMD_UPDATE_ROLE_ADDRESS_NTF, ntf)
		connectResp.Success = true
	}
	resp.Buffer, _ = proto.Marshal(connectResp)
}

func (g *Gate) handleMessage(role *session, message *msg.ClientToGate) {
	if message.Command == msg.CMD_ROLE_ENTER {
		roleEnterResp := &msg.RoleEnterResp{}
		resp := &msg.GateToClient{}
		resp.Command = msg.CMD_RESP
		resp.Sequence = message.Sequence
		defer g.sendSessionResp(role, resp)
		roleEnterReq := &msg.RoleEnterReq{}
		if err := proto.Unmarshal(message.Request, roleEnterReq); err != nil {
			logger.Info(err)
			return
		}
		loadRoleReq := &msg.LoadRoleReq{
			RoleUUID: roleEnterReq.RoleUUID,
		}
		loadRoleResp := &msg.LoadRoleResp{}
		if err := node.Instance.Net.Request(msg.CMD_GD_LOAD_ROLE_REQ, loadRoleReq, loadRoleResp); err != nil {
			logger.Info(err)
			return
		}
		roleEnterResp.Success = loadRoleResp.Success
		roleEnterResp.Role = loadRoleResp.Role
		if loadRoleResp.Success {
			getMapIDReq := &msg.GetMapIDReq{
				MapUUID: loadRoleResp.MapUUID,
			}
			getMapIDResp := &msg.GetMapIDResp{}
			err := node.Instance.Net.RequestServer(msg.CMD_GET_MAP_ID, getMapIDReq, getMapIDResp, loadRoleResp.WorldID)
			if err != nil {
				logger.Info(err)
				return
			}
			roleEnterResp.MapID = getMapIDResp.MapID
		}
		resp.Buffer, _ = proto.Marshal(roleEnterResp)
		if loadRoleResp.Success == false {
			logger.Info("Role Enter Fail! UUID:" + role.info.RoleUUID)
		} else {
			role.info.Address.Entity = loadRoleResp.WorldID
			ntf := &msg.UpdateRoleAddressNTF{
				RoleUUID: roleEnterReq.RoleUUID,
				Address:  role.info.Address,
			}
			node.Instance.Net.Notify(msg.CMD_UPDATE_ROLE_ADDRESS_NTF, ntf)
		}
	} else if message.Request != nil {
		buffer, err := node.Instance.Net.RequestClientBytesToServer(message.Command, role.info.RoleUUID, message.Request, role.info.Address.Entity)
		if err != nil {
			logger.Debug(err)
		} else {
			resp := &msg.GateToClient{}
			resp.Command = msg.CMD_RESP
			resp.Sequence = message.Sequence
			resp.Buffer = buffer
			respBuffer, err := proto.Marshal(resp)
			if err != nil {
				logger.Debug(err)
			} else {
				if role.conn != nil {
					innet.SendBytesHelper(role.conn, respBuffer)
				} else if role.webconn != nil {
					innet.SendWebBytesHelper(role.webconn, respBuffer)
				}
			}
		}
	} else if message.Notify != nil {
		err := node.Instance.Net.NotifyClientBytesToServer(message.Command, role.info.RoleUUID, message.Notify, role.info.Address.Entity)
		if err != nil {
			logger.Debug(err)
		}
	}
}

func (g *Gate) pushNewSessionCert(role *session) {
	buff, _ := proto.Marshal(&msg.SessionCert{
		UUID: role.info.RoleUUID,
		Key:  role.cert.Key,
	})
	message := &msg.Message{
		Header: &msg.MessageHeader{
			Command:  msg.CMD_SESSION_CERT_NTF,
			Sequence: 0,
			From:     global.CurrentServerID,
		},
		Buffer: buff,
	}
	g.sendSessionResp(role, message)
}

func (g *Gate) sendSessionResp(s *session, resp proto.Message) error {
	if s.conn != nil {
		return g.sendResp(s.conn, resp)
	} else if s.webconn != nil {
		return g.sendWebResp(s.webconn, resp)
	}
	return errors.New("NO CONN")
}

func (g *Gate) sendResp(conn *net.TCPConn, resp proto.Message) error {
	buff, _ := proto.Marshal(resp)
	return innet.SendBytesHelper(conn, buff)
}

func (g *Gate) sendWebResp(conn *websocket.Conn, resp proto.Message) error {
	buff, _ := proto.Marshal(resp)
	return innet.SendWebBytesHelper(conn, buff)
}
