package network

import (
	"fmt"
	"net"
	"ssemu/internal/utils"
	"strconv"
	"strings"
)

type Endpoint struct {
	address uint32
	port    uint16
}

func NewEndpoint(address uint32, port uint16) Endpoint {
	return Endpoint{
		address: address,
		port:    port,
	}
}

func NewEndpointFromString(ipv4 string, port uint16) Endpoint {
	return Endpoint{
		address: utils.IPToU32(ipv4),
		port:    uint16(port),
	}
}

func NewEndpointFromAddress(addr net.Addr) Endpoint {
	endpointStr := addr.String()
	index := strings.Index(endpointStr, ":")
	ipStr := endpointStr[0:index]
	ipU32 := utils.IPToU32(ipStr)
	port, _ := strconv.ParseUint(endpointStr[index+1:], 10, 16)
	return Endpoint{
		address: ipU32,
		port:    uint16(port),
	}
}

func (e Endpoint) GetString() string {
	return fmt.Sprintf("%s:%d", utils.IPToString(e.address), e.port)
}

func (e Endpoint) GetAddress() uint32 { return e.address }

func (e Endpoint) GetPort() uint16 { return e.port }
