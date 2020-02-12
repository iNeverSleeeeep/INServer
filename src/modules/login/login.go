package login

import (
	"INServer/src/common/dbobj"
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/protect"
	"INServer/src/common/util"
	"INServer/src/common/uuid"
	"INServer/src/dao"
	"INServer/src/modules/cluster"
	"INServer/src/modules/innet"
	"INServer/src/modules/node"
	"INServer/src/proto/data"
	"INServer/src/proto/db"
	"INServer/src/proto/msg"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

var Instance *Login
var upgrader = websocket.Upgrader{}

type (
	Login struct {
		listener *net.TCPListener
		DB       *dbobj.DBObject
	}
)

func New() *Login {
	l := new(Login)
	l.DB = dbobj.New()
	l.DB.Open(global.CurrentServerConfig.LoginConfig.Database, global.DatabaseSchema)
	return l
}

func (l *Login) Start() {
	// TCP
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: int(global.CurrentServerConfig.LoginConfig.Port)})
	if err != nil {
		log.Fatalln(err)
	}
	l.listener = listener

	logger.Info("登录服务器 启动 监听端口:" + strconv.Itoa(int(global.CurrentServerConfig.LoginConfig.Port)))
	go func() {
		for {
			conn, err := l.listener.AcceptTCP()
			if err != nil {
				logger.Debug("listener.Accept 连接错误")
			} else {
				go l.handleConnect(conn)
			}
		}
	}()

	// WebSocket
	if global.CurrentServerConfig.LoginConfig.WebPort > 0 {
		http.HandleFunc("/", l.handleWebConnect)
		go http.ListenAndServe(fmt.Sprintf(":%d", global.CurrentServerConfig.LoginConfig.WebPort), nil)
		logger.Info("登录服务器 监听端口:" + strconv.Itoa(int(global.CurrentServerConfig.LoginConfig.WebPort)))
	}
}

func (l *Login) handleWebConnect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("upgrade:", err)
		return
	}
	defer c.Close()
	defer protect.CatchPanic()
	var player **data.Player
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
		message := &msg.ClientToLogin{}
		err := proto.Unmarshal(buf[2:size+2], message)
		if err != nil {
			logger.Debug("proto解析失败")
			return
		}
		l.handleWebMessage(c, message, player)

		copy(buf[0:], buf[size+2:current])
		current = current - int(size) - 2
	}
}

func (l *Login) handleConnect(conn *net.TCPConn) {
	defer conn.Close()
	defer protect.CatchPanic()
	var player **data.Player
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

		message := &msg.ClientToLogin{}
		err := proto.Unmarshal(buf[2:size+2], message)
		if err != nil {
			logger.Debug("proto解析失败")
			return
		}
		l.handleMessage(conn, message, player)

		copy(buf[0:], buf[size+2:current])
		current = current - int(size) - 2
	}
}

func newAccount(name string, passwordHash string) *db.DBAccount {
	account := new(db.DBAccount)
	account.Name = name
	account.PasswordHash = passwordHash
	account.PlayerUUID = uuid.New()
	return account
}

func (l *Login) handleMessageImpl(message *msg.ClientToLogin, resp *msg.LoginToClient, player **data.Player) {
	success := false
	var account *db.DBAccount
	var err error
	if message.Logon != nil {
		account = newAccount(message.Logon.Name, message.Logon.PasswordHash)
		createPlayerReq := &msg.CreatePlayerReq{
			PlayerUUID: account.PlayerUUID,
		}
		createPlayerResp := &msg.CreatePlayerResp{}
		err = node.Instance.Net.Request(msg.CMD_LD_CREATE_PLAYER_REQ, createPlayerReq, createPlayerResp)
		if err != nil {
			logger.Debug(err)
		} else if createPlayerResp.Success {
			err = dao.AccountInsert(l.DB, account)
			if err != nil {
				logger.Debug(err)
			} else {
				success = true
			}
		}
	} else if message.Login != nil {
		account, err = dao.AccountQuery(l.DB, message.Login.Name)
		if err != nil {
			logger.Debug(err)
		} else if account != nil {
			success = account.PasswordHash == message.Login.PasswordHash
			if success == false {
				logger.Info(fmt.Sprintf("密码错误 %s %s", account.PasswordHash, message.Login.PasswordHash))
			} else {
				// 推送消息到其他服务器
				logger.Info("客户端登录成功:" + account.PlayerUUID)

				loadPlayerReq := &msg.LoadPlayerReq{
					PlayerUUID: account.PlayerUUID,
				}
				loadPlayerResp := &msg.LoadPlayerResp{}
				err := node.Instance.Net.Request(msg.CMD_GD_LOAD_PLAYER_REQ, loadPlayerReq, loadPlayerResp)
				if err != nil {
					logger.Info(err)
					success = false
				} else if loadPlayerResp.Success == false {
					success = false
				} else {
					*player = loadPlayerResp.Player
				}
			}
		}
	} else if message.ChangePassword != nil {
		account, err := dao.AccountQuery(l.DB, message.Login.Name)
		if err != nil {
			logger.Debug(err)
		} else if account != nil && account.PasswordHash == message.ChangePassword.OldPasswordHash {
			account.PasswordHash = message.ChangePassword.NewPasswordHash
			err = dao.AccountUpdate(l.DB, account)
			if err != nil {
				logger.Debug(err)
			} else {
				success = true
			}
		}
	} else if message.CreateRole != nil {
		if *player != nil {
			createRoleResp := &msg.CreateRoleResp{}
			createRoleReq := &msg.CreateRoleReq{}
			err := node.Instance.Net.Request(msg.CMD_GD_CREATE_ROLE_REQ, createRoleReq, createRoleResp)
			if err != nil {
				logger.Info(err)
			} else if createRoleResp.Success {
				(*player).RoleList = append((*player).RoleList, createRoleResp.Role)
			}
		}
	} else if message.EnterGame != nil {
		if *player != nil {
			invalidRole := true
			for _, role := range (*player).RoleList {
				if role.RoleUUID == message.EnterGame.RoleUUID {
					invalidRole = false
					break
				}
			}
			if invalidRole == false {
				gateID := l.selectGate()
				if gateID != global.InvalidServerID {
					cert := &msg.SessionCert{
						UUID: message.EnterGame.RoleUUID,
						Key:  util.GetRandomString(global.CERT_KEY_LEN),
					}
					resp.SessionCert = cert

					ip, port, webport := cluster.GetGatePublicAddress(gateID)
					resp.GateIP, resp.GatePort, resp.GateWebPort = ip, int32(port), int32(webport)
					message := &msg.LoginToGate{Cert: cert}
					node.Instance.Net.NotifyServer(msg.CMD_SESSION_CERT_NTF, message, gateID)
					logger.Info(fmt.Sprintf("玩家登录成功 UUID:%s Name:%s CertKey:%s", cert.UUID, account.Name, cert.Key))
				} else {
					logger.Error("没有找到门服务器")
					success = false
				}
			}
		}
	}

	resp.Player = *player
	resp.Success = success
}

func (l *Login) handleWebMessage(conn *websocket.Conn, message *msg.ClientToLogin, player **data.Player) {
	resp := &msg.LoginToClient{}
	defer l.sendWebResponce(conn, resp)
	defer protect.CatchPanic()
	l.handleMessageImpl(message, resp, player)
}

func (l *Login) handleMessage(conn *net.TCPConn, message *msg.ClientToLogin, player **data.Player) {
	resp := &msg.LoginToClient{}
	defer l.sendResponce(conn, resp)
	defer protect.CatchPanic()
	l.handleMessageImpl(message, resp, player)
}

func (l *Login) sendWebResponce(conn *websocket.Conn, message *msg.LoginToClient) {
	buff, _ := proto.Marshal(message)
	buffer, _ := proto.Marshal(&msg.Message{
		Header: &msg.MessageHeader{
			Command:  msg.CMD_SESSION_CERT_NTF,
			Sequence: 0,
			From:     global.CurrentServerID,
		},
		Buffer: buff,
	})
	innet.SendWebBytesHelper(conn, buffer)
}

func (l *Login) sendResponce(conn *net.TCPConn, message *msg.LoginToClient) {
	buff, _ := proto.Marshal(message)
	buffer, _ := proto.Marshal(&msg.Message{
		Header: &msg.MessageHeader{
			Command:  msg.CMD_SESSION_CERT_NTF,
			Sequence: 0,
			From:     global.CurrentServerID,
		},
		Buffer: buff,
	})
	innet.SendBytesHelper(conn, buffer)
}

func (l *Login) selectGate() int32 {
	gates := cluster.RunningGates()
	if len(gates) > 0 {
		index := time.Now().UnixNano() % int64(len(gates))
		return gates[index]
	} else {
		return global.InvalidServerID
	}
}
