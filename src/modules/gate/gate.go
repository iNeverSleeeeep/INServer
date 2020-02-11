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
		players  map[string]*session
		serv     *services
	}
)

func New() *Gate {
	g := new(Gate)
	g.players = make(map[string]*session)
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
	node.Instance.Net.Listen(msg.CMD_UPDATE_PLAYER_ADDRESS_NTF, g.onUpdatePlayerAddressNTF)
}

func (g *Gate) onUpdatePlayerAddressNTF(header *msg.MessageHeader, buffer []byte) {
	ntf := &msg.UpdatePlayerAddressNTF{}
	err := proto.Unmarshal(buffer, ntf)
	if err != nil {
		logger.Error(err)
		return
	}

	if session, ok := g.players[ntf.PlayerUUID]; ok {
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
			if player, ok := g.players[*uuid]; ok {
				if player.info.State == data.SessionState_Online {
					g.handleMessage(player, message)
				} else {
					logger.Debug("客户端状态错误:" + c.RemoteAddr().String() + " 当前状态:" + player.info.State.String())
					return
				}
			}
		} else {
			logger.Debug("客户端需要先发送Connect协议:" + c.RemoteAddr().String())
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
			if player, ok := g.players[*uuid]; ok {
				if player.info.State == data.SessionState_Online {
					g.handleMessage(player, message)
				} else {
					logger.Debug("客户端状态错误:" + conn.RemoteAddr().String() + " 当前状态:" + player.info.State.String())
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
		if player, ok := g.players[*uuid]; ok {
			player.conn = nil
			player.webconn = nil
			player.info.State = data.SessionState_OutOfContact
			player.cert.OutOfDateTime = generateCertOutOfDateTime()
		}
	}
}

func (g *Gate) tickOutOfContact() {
	go func() {
		for {
			time.Sleep(time.Second)
			now := time.Now().Unix()
			var outOfDatePlayers []string
			for uuid, player := range g.players {
				if player.info.State == data.SessionState_OutOfContact && player.cert.OutOfDateTime < now {
					outOfDatePlayers = append(outOfDatePlayers, uuid)
				}
			}
			for _, uuid := range outOfDatePlayers {
				delete(g.players, uuid)
			}
		}
	}()
}

func (g *Gate) onPlayerConnect(player *session) (*data.Player, error) {
	// 推送消息到其他服务器
	logger.Info("客户端登录成功:" + player.info.UUID)

	loadPlayerReq := &msg.LoadPlayerReq{
		PlayerUUID: player.info.UUID,
	}
	loadPlayerResp := &msg.LoadPlayerResp{}
	err := node.Instance.Net.Request(msg.CMD_GD_LOAD_PLAYER_REQ, loadPlayerReq, loadPlayerResp)
	if err != nil {
		return nil, err
	}

	return loadPlayerResp.Player, nil
}

func (g *Gate) onPlayerReconnect(player *session) error {
	// 推送从消息到其他服务器
	return nil
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
	player, ok := g.players[*uuid]
	if ok == false {
		addr := ""
		if conn != nil {
			addr = conn.RemoteAddr().String()
		} else if webconn != nil {
			addr = webconn.RemoteAddr().String()
		}
		logger.Debug("拒绝连接，没有数据:" + addr + " uuid:" + *uuid)
		return
	} else if message.SessionCert.Key != player.cert.Key {
		addr := ""
		if conn != nil {
			addr = conn.RemoteAddr().String()
		} else if webconn != nil {
			addr = webconn.RemoteAddr().String()
		}
		logger.Info(fmt.Sprintf("客户端秘钥错误 addr:%d uuid:%s client:%s server:%s", addr, *uuid, message.SessionCert.Key, player.cert.Key))
		return
	} else {
		player.conn = conn
		player.webconn = webconn
		player.generateNewCertKey()
		oldState := player.info.State
		player.info.State = data.SessionState_Online
		player.info.UUID = *uuid
		player.info.Address = &data.PlayerAddress{
			Gate:   global.CurrentServerID,
			Entity: global.InvalidServerID,
		}
		ntf := &msg.UpdatePlayerAddressNTF{
			PlayerUUID: *uuid,
			Address:    player.info.Address,
		}
		node.Instance.Net.Notify(msg.CMD_UPDATE_PLAYER_ADDRESS_NTF, ntf)
		if oldState == data.SessionState_Offline || oldState == data.SessionState_Connected || oldState == data.SessionState_Online {
			playerData, err := g.onPlayerConnect(player)
			if err != nil {
				logger.Debug(err)
				return
			}
			connectResp.Success = true
			connectResp.Player = playerData
		} else if oldState == data.SessionState_OutOfContact {
			err := g.onPlayerReconnect(player)
			if err != nil {
				logger.Debug(err)
				return
			}
			connectResp.Success = true
		}
	}
	resp.Buffer, _ = proto.Marshal(connectResp)
}

func (g *Gate) handleMessage(player *session, message *msg.ClientToGate) {
	if message.Command == msg.CMD_ROLE_ENTER {
		roleEnterResp := &msg.RoleEnterResp{}
		resp := &msg.GateToClient{}
		resp.Command = msg.CMD_RESP
		resp.Sequence = message.Sequence
		defer g.sendSessionResp(player, resp)
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
			logger.Info("Role Enter Fail! UUID:" + player.info.UUID)
		} else {
			player.info.Address.Entity = loadRoleResp.WorldID
			ntf := &msg.UpdatePlayerAddressNTF{
				PlayerUUID: player.info.UUID,
				RoleUUID:   roleEnterReq.RoleUUID,
				Address:    player.info.Address,
			}
			node.Instance.Net.Notify(msg.CMD_UPDATE_PLAYER_ADDRESS_NTF, ntf)
		}
	} else if message.Request != nil {
		buffer, err := node.Instance.Net.RequestClientBytesToServer(message.Command, player.info.UUID, message.Request, player.info.Address.Entity)
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
				if player.conn != nil {
					innet.SendBytesHelper(player.conn, respBuffer)
				} else if player.webconn != nil {
					innet.SendWebBytesHelper(player.webconn, respBuffer)
				}
			}
		}
	} else if message.Notify != nil {
		err := node.Instance.Net.NotifyClientBytesToServer(message.Command, player.info.UUID, message.Notify, player.info.Address.Entity)
		if err != nil {
			logger.Debug(err)
		}
	}
}

func (g *Gate) pushNewSessionCert(player *session) {
	buff, _ := proto.Marshal(&msg.SessionCert{
		UUID: player.info.UUID,
		Key:  player.cert.Key,
	})
	message := &msg.Message{
		Header: &msg.MessageHeader{
			Command:  msg.CMD_SESSION_CERT_NTF,
			Sequence: 0,
			From:     global.CurrentServerID,
		},
		Buffer: buff,
	}
	g.sendSessionResp(player, message)
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
