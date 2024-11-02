package nat

import (
	"fmt"
	"ssemu/internal/domain"
)

type NatInfo struct {
	UserId         domain.UserId
	PublicAddress  uint32
	PublicPort     uint16
	PrivateAddress uint32
	PrivatePort    uint16
	NatUnk         uint16
	ConnectionType byte
}

var infoMap = make(map[domain.UserId]NatInfo)

func Set(info NatInfo) {
	infoMap[info.UserId] = info
}

func Get(userId domain.UserId) (NatInfo, error) {
	info, found := infoMap[userId]
	if !found {
		return NatInfo{}, fmt.Errorf("user %d nat info not found", userId)
	}
	return info, nil
}
