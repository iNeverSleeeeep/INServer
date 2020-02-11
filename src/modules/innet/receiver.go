package innet

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/proto/msg"
	"encoding/binary"
	"net"

	"github.com/gogo/protobuf/proto"
)

type (
	serverPackageCache struct {
		currentUniqueID uint64
		packages        map[uint64]map[int32]*msg.Package
	}
	receiver struct {
		innet        *INNet
		packageCache map[int32]*serverPackageCache
	}
)

func newReceiver(innet *INNet) *receiver {
	r := new(receiver)
	r.innet = innet
	r.packageCache = make(map[int32]*serverPackageCache)
	return r
}

func (r *receiver) start() {
	go r.receiveLoop()
}

func (r *receiver) receiveLoop() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: int(global.CurrentServerID) + recvport})
	if err != nil {
		logger.Fatal(err)
	}
	var buf = make([]byte, 65536)
	current := 0
	var udpaddr *net.UDPAddr
	for {
		// 等待读取数据长度
		for current < 2 {
			n, addr, _ := conn.ReadFromUDP(buf[current:])
			udpaddr = addr
			current = current + n
		}

		// 等待读取数据
		size := binary.BigEndian.Uint16(buf[:2])
		for (current - 2) < int(size) {
			n, _ := conn.Read(buf[current:])
			current = current + n
		}

		pkg := &msg.Package{}
		proto.Unmarshal(buf[2:size+2], pkg)
		r.innet.retry.handlePackageReceive(pkg)
		if pkg.Buffer != nil {
			port := udpaddr.Port + recvport - sendport
			addr := &net.UDPAddr{IP: udpaddr.IP, Port: port}
			r.innet.sender.ack(addr, pkg)
			r.handlePackage(pkg)
		}

		copy(buf[0:], buf[size+2:current])
		current = current - int(size) - 2
	}
}

func (r *receiver) handlePackage(pkg *msg.Package) {
	if _, ok := r.packageCache[pkg.From]; ok == false {
		r.packageCache[pkg.From] = new(serverPackageCache)
		r.packageCache[pkg.From].packages = make(map[uint64]map[int32]*msg.Package)
		r.packageCache[pkg.From].currentUniqueID = 1
	}
	packages := r.packageCache[pkg.From]

	// 如果收到的UniqueID为0直接处理，因为这是服务器启动后的第一个包
	if pkg.UniqueID == 0 {
		msg := &msg.Message{}
		proto.Unmarshal(pkg.Buffer, msg)
		r.innet.handleMessage(msg)
	} else if pkg.UniqueID >= packages.currentUniqueID {
		if _, ok := packages.packages[pkg.UniqueID]; ok == false {
			packages.packages[pkg.UniqueID] = make(map[int32]*msg.Package)
		}
		sequencePackages := packages.packages[pkg.UniqueID]
		if _, ok := sequencePackages[pkg.Index]; ok == false {
			sequencePackages[pkg.Index] = pkg
			if len(sequencePackages) == int(pkg.Total) && pkg.UniqueID == packages.currentUniqueID {
				r.onPackagesReceiveFinished(sequencePackages)
				delete(packages.packages, packages.currentUniqueID)
				packages.currentUniqueID++
				r.testCurrentPackages(packages)
			}
		}
	}
}

func (r *receiver) resetServer(serverID int32) {
	if _, ok := r.packageCache[serverID]; ok {
		delete(r.packageCache, serverID)
	}
}

func (r *receiver) testCurrentPackages(serverPackages *serverPackageCache) {
	if packages, ok := serverPackages.packages[serverPackages.currentUniqueID]; ok {
		for _, packageCache := range packages {
			if packageCache.Total == int32(len(packages)) {
				r.onPackagesReceiveFinished(packages)
				delete(serverPackages.packages, serverPackages.currentUniqueID)
				serverPackages.currentUniqueID++
				r.testCurrentPackages(serverPackages)
			}
			break
		}
	}
}

func (r *receiver) onPackagesReceiveFinished(packages map[int32]*msg.Package) {
	var buffer []byte
	for index := int32(0); index < int32(len(packages)); index++ {
		buffer = append(buffer, packages[index].Buffer...)
	}
	msg := &msg.Message{}
	proto.Unmarshal(buffer, msg)
	r.innet.handleMessage(msg)
}
