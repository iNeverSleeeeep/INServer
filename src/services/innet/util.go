package innet

import (
	"INServer/src/common/logger"
	"errors"
	"net"
)

const (
	timeout   int64 = 10 * 1E6 // 纳秒
	sliceSize int   = 1000     // MTU 1472

	sendport   = 11000
	recvport   = 12000
	expvarport = 13000
)

func getIP() []byte {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		logger.Fatal(err)
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.To4()
					}
				}
			}
		}
	}

	logger.Fatal(errors.New("No Valid IPV4!"))
	return nil
}
