// 参考了github.com/satori/go.uuid v1的代码，进行了一些修改减少了位数，可以在服务器群中生成一个唯一的UUID
package uuid

import (
	"INServer/src/common/global"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

const epochStart = 122192928000000000

type epochFunc func() time.Time
type hwAddrFunc func() (net.HardwareAddr, error)

type (
	rfc4122Generator struct {
		clockSequenceOnce sync.Once
		hardwareAddrOnce  sync.Once
		storageMutex      sync.Mutex

		rand io.Reader

		epochFunc     epochFunc
		hwAddrFunc    hwAddrFunc
		lastTime      uint64
		clockSequence uint16
		hardwareAddr  [6]byte
	}
)

var generator = newRFC4122Generator()

func newRFC4122Generator() *rfc4122Generator {
	return &rfc4122Generator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       rand.Reader,
	}
}

func New() string {
	u := make([]byte, 8)
	timeNow, clockSeq, _ := generator.getClockSequence()
	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[4:], clockSeq)
	binary.BigEndian.PutUint16(u[6:], uint16(global.ServerID))

	return strconv.FormatUint(binary.BigEndian.Uint64(u), 36)
}

func (g *rfc4122Generator) getClockSequence() (uint64, uint16, error) {
	var err error
	g.clockSequenceOnce.Do(func() {
		buf := make([]byte, 2)
		if _, err = io.ReadFull(g.rand, buf); err != nil {
			return
		}
		g.clockSequence = binary.BigEndian.Uint16(buf)
	})
	if err != nil {
		return 0, 0, err
	}

	g.storageMutex.Lock()
	defer g.storageMutex.Unlock()

	timeNow := g.getEpoch()
	// Clock didn't change since last UUID generation.
	// Should increase clock sequence.
	if timeNow <= g.lastTime {
		g.clockSequence++
	}
	g.lastTime = timeNow

	return timeNow, g.clockSequence, nil
}

func (g *rfc4122Generator) getEpoch() uint64 {
	return epochStart + uint64(g.epochFunc().UnixNano()/100)
}

func defaultHWAddrFunc() (net.HardwareAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return []byte{}, err
	}
	for _, iface := range ifaces {
		if len(iface.HardwareAddr) >= 6 {
			return iface.HardwareAddr, nil
		}
	}
	return []byte{}, errors.New("uuid: no HW address found")
}
