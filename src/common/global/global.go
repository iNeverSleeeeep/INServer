package global

import (
	"INServer/src/proto/etc"
)

const (
	CenterID        int32  = 0
	InvalidServerID int32  = -1
	InvalidServer   string = "InvalidServer"

	CenterServer   string = "Center"
	GateServer     string = "Gate"
	LoginServer    string = "Login"
	ChatServer     string = "Chat"
	GPSServer      string = "GPS"
	WorldServer    string = "World"
	DatabaseServer string = "Database"
	WebServer      string = "Web"

	ZoneStateOpen   string = "Open"
	ZoneStateClosed string = "Closed"

	DatabaseSchema string = "indb"

	SERVER_ID_MAX = 999
	CERT_KEY_LEN  = 10
)

var (
	ServerID     int32  = 0
	ServerType   string = InvalidServer
	ServerConfig *etc.ServerConfig
	Servers      []*etc.Server
	Zones        []*etc.Zone

	CenterIP string = "127.0.0.1"
)
