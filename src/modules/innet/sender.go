package innet

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/proto/msg"
	"encoding/binary"
	"errors"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type (
	sender struct {
		innet *INNet
	}
)

func newSender(innet *INNet) *sender {
	s := new(sender)
	s.innet = innet
	return s
}

func SendUDPBytesHelper(addr *net.UDPAddr, bytes []byte) error {
	sizebuf := make([]byte, 2)
	binary.BigEndian.PutUint16(sizebuf, uint16(len(bytes)))
	_, err := udpconn.WriteToUDP(append(sizebuf, bytes...), addr)
	if err != nil {
		logger.Debug(err)
	}
	return err
}

func SendBytesHelper(conn net.Conn, bytes []byte) error {
	sizebuf := make([]byte, 2)
	binary.BigEndian.PutUint16(sizebuf, uint16(len(bytes)))
	_, err := conn.Write(append(sizebuf, bytes...))
	if err != nil {
		logger.Debug(err)
	}
	return err
}

func SendWebBytesHelper(conn *websocket.Conn, bytes []byte) error {
	sizebuf := make([]byte, 2)
	binary.BigEndian.PutUint16(sizebuf, uint16(len(bytes)))
	err := conn.WriteMessage(websocket.BinaryMessage, append(sizebuf, bytes...))
	if err != nil {
		logger.Debug(err)
	}
	return err
}

func checkServerValid(svr *server) bool {
	return svr.info.State != msg.ServerState_Offline
}

// 这个是发给center用的
func (s *sender) sendstate(svr *server, buffer []byte) error {
	var packages []*msg.Package
	packages = append(packages, &msg.Package{
		UniqueID: 0,
		From:     int32(global.ServerID),
		Index:    0,
		Total:    1,
		Buffer:   buffer,
	})
	s.sendPackages(svr.addr, packages)
	return nil
}

func (s *sender) send(svr *server, buffer []byte) error {
	if checkServerValid(svr) == false {
		return errors.New("Server Offline")
	}
	svr.packageID++
	slices := cutslices(buffer)
	var packages []*msg.Package
	total := int32(len(slices))
	for index, slice := range slices {
		packages = append(packages, &msg.Package{
			UniqueID: svr.packageID,
			From:     int32(global.ServerID),
			Index:    int32(index),
			Total:    total,
			Buffer:   slice,
		})
	}
	s.sendPackages(svr.addr, packages)

	s.innet.retry.addPackagesCache(svr.packageID, &packageCache{
		addr:       svr.addr,
		packages:   packages,
		retryTime:  time.Now().UnixNano() + timeout,
		toServerID: svr.info.ServerID,
	})
	return nil
}

func (s *sender) ack(addr *net.UDPAddr, pkg *msg.Package) {
	buffer, _ := proto.Marshal(&msg.Package{
		UniqueID: pkg.UniqueID,
		From:     int32(global.ServerID),
		Index:    pkg.Index,
		Total:    pkg.Total,
	})
	SendUDPBytesHelper(addr, buffer)
}

func (s *sender) sendPackages(addr *net.UDPAddr, packages []*msg.Package) {
	for _, pkg := range packages {
		buf, _ := proto.Marshal(pkg)
		SendUDPBytesHelper(addr, buf)
	}
}

func cutslices(buffer []byte) [][]byte {
	var slices [][]byte
	from := 0
	for len(buffer[from:]) > sliceSize {
		to := from + sliceSize
		slices = append(slices, buffer[from:to])
		from = to
	}
	slices = append(slices, buffer[from:])
	return slices
}
