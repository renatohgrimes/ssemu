package utils

import (
	"encoding/binary"
	"fmt"
	"net"
)

func IPToU32(ipStr string) uint32 {
	return binary.LittleEndian.Uint32(net.ParseIP(ipStr).To4())
}

func IPToString(address uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(address),
		byte(address>>8),
		byte(address>>16),
		byte(address>>24))
}
