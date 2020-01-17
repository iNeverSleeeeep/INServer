package innet

import (
	"INServer/src/proto/msg"
	"net"
	"time"
)

type (
	packageCache struct {
		addr       *net.UDPAddr
		packages   []*msg.Package
		retryTime  int64
		toServerID int32
	}

	retry struct {
		innet          *INNet
		cachedPackages map[uint64]*packageCache
		stoped         bool
	}
)

func newRetry(innet *INNet) *retry {
	r := new(retry)
	r.innet = innet
	r.cachedPackages = make(map[uint64]*packageCache)
	r.startTickRetry()
	return r
}

func (r *retry) addPackagesCache(packageID uint64, cache *packageCache) {
	r.cachedPackages[packageID] = cache
}

func (r *retry) resetServer(serverID int32) {
	packageIDList := make([]uint64, 0)
	for packageID, cache := range r.cachedPackages {
		if cache.toServerID == serverID {
			packageIDList = append(packageIDList, packageID)
		}
	}
	for _, packageID := range packageIDList {
		delete(r.cachedPackages, packageID)
	}
}

func (r *retry) handlePackageReceive(receivedPkg *msg.Package) {
	if cache, ok := r.cachedPackages[receivedPkg.UniqueID]; ok {
		for index, pkg := range cache.packages {
			if pkg.UniqueID == receivedPkg.UniqueID {
				cache.packages = append(cache.packages[:index], cache.packages[index+1:]...)
				if len(cache.packages) == 0 {
					delete(r.cachedPackages, receivedPkg.UniqueID)
				}
				break
			}
		}
	}
}

func (r *retry) startTickRetry() {
	r.stoped = false
	go func() {
		for r.stoped == false {
			time.Sleep(time.Microsecond)
			current := time.Now().UnixNano()
			for _, cache := range r.cachedPackages {
				retryTime := cache.retryTime
				if retryTime < current {
					r.innet.sender.sendPackages(cache.addr, cache.packages)
					cache.retryTime = current + timeout
				}
			}
		}
	}()
}
