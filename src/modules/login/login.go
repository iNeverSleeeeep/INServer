package login

import (
	"INServer/src/common/dbobj"
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/common/protect"
	"INServer/src/common/util"
	"INServer/src/common/uuid"
	"INServer/src/dao"
	"INServer/src/modules/innet"
	"INServer/src/modules/node"
	"INServer/src/proto/db"
	"INServer/src/proto/msg"
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
)

var Instance *Login

type (
	Login struct {
		listener *net.TCPListener
		DB       *dbobj.DBObject
	}
)

func New() *Login {
	l := new(Login)
	l.DB = dbobj.New()
	l.DB.Open(global.ServerConfig.LoginConfig.Database, global.DatabaseSchema)
	return l
}

func (l *Login) Start() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: int(global.ServerConfig.LoginConfig.Port)})
	if err != nil {
		log.Fatalln(err)
	}
	l.listener = listener

	logger.Info("登录服务器 启动 监听端口:" + strconv.Itoa(int(global.ServerConfig.LoginConfig.Port)))
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
}

func (l *Login) handleConnect(conn *net.TCPConn) {
	defer conn.Close()
	defer protect.CatchPanic()
	var buf = make([]byte, 65536)
	current := 0
	for {
		// 等待读取数据长度
		for current < 2 {
			n, err := conn.Read(buf[current:])
			if err != nil {
				logger.Debug(err)
				return
			}
			current = current + n
		}

		// 等待读取数据
		size := binary.BigEndian.Uint16(buf[:2])
		for (current - 2) < int(size) {
			n, err := conn.Read(buf[current:])
			if err != nil {
				logger.Debug(err)
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
		l.handleMessage(conn, message)

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

func (l *Login) handleMessage(conn *net.TCPConn, message *msg.ClientToLogin) {
	resp := &msg.LoginToClient{}
	defer l.sendResponce(conn, resp)
	defer protect.CatchPanic()

	success := false
	var account *db.DBAccount
	var err error
	if message.Logon != nil {
		account = newAccount(message.Logon.Name, message.Logon.PasswordHash)
		createPlayerReq := &msg.CreatePlayerReq{
			PlayerUUID: account.PlayerUUID,
		}
		createPlayerResp := &msg.CreatePlayerResp{}
		err = node.Instance.Net.Request(msg.Command_LD_CREATE_PLAYER_REQ, createPlayerReq, createPlayerResp)
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
	}

	if success {
		if message.Login != nil {
			gateID := l.selectGate()
			if gateID != global.InvalidServerID {
				cert := &msg.SessionCert{
					UUID: account.PlayerUUID,
					Key:  util.GetRandomString(global.CERT_KEY_LEN),
				}
				resp.SessionCert = cert

				ip, port := node.Instance.Net.GetGatePublicAddress(gateID)
				resp.GateIP, resp.GatePort = ip, int32(port)
				message := &msg.LoginToGate{
					Cert: cert,
				}
				node.Instance.Net.NotifyServer(msg.Command_SESSION_CERT_NTF, message, gateID)
			} else {
				success = false
			}
		}
		resp.Success = success
	}
}

func (l *Login) sendResponce(conn *net.TCPConn, message *msg.LoginToClient) {
	buff, _ := proto.Marshal(message)
	buffer, _ := proto.Marshal(&msg.Message{
		Header: &msg.MessageHeader{
			Command:  msg.Command_SESSION_CERT_NTF,
			Sequence: 0,
			From:     global.ServerID,
		},
		Buffer: buff,
	})
	innet.SendBytesHelper(conn, buffer)
}

func (l *Login) selectGate() int32 {
	gates := node.Instance.Net.Gates()
	if len(gates) > 0 {
		index := time.Now().UnixNano() % int64(len(gates))
		return gates[index]
	} else {
		return global.InvalidServerID
	}
}
