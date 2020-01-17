package unittest

import (
	"INServer/src/common/logger"
	"net"
	"testing"
)

func newConn() *net.UDPConn {
	raddr, err := net.ResolveUDPAddr("udp4", ":10000")
	if err != nil {
		logger.Debug(err)
		return nil
	}
	//laddr, err := net.ResolveUDPAddr("udp4", ":10001")
	//if err != nil {
	//	logger.Debug(err)
	//  return nil
	//}
	conn, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		logger.Debug(err)
		return nil
	}
	return conn
}

func TestDialUDP(t *testing.T) {
	conn := newConn()
	if conn == nil {
		t.FailNow()
	}
	conn = newConn()
	if conn == nil {
		t.FailNow()
	}
	conn = newConn()
	if conn == nil {
		t.FailNow()
	}
	conn = newConn()
	if conn == nil {
		t.FailNow()
	}
}
