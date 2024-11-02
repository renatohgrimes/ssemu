package handlers

import (
	"context"
	"ssemu/internal/app/game/nat"
	"ssemu/internal/network"
)

type NatInfoReqHandler struct{}

const NatInfoReq network.PacketId = 0x0F

func (h *NatInfoReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	privateAddress := p.ReadU32()
	privatePort := p.ReadU16()
	publicAddress := p.ReadU32()
	publicPort := p.ReadU16()
	unk := p.ReadU16()
	connectionType := p.ReadU8()

	privateEndpoint := network.NewEndpoint(privateAddress, privatePort)
	publicEndpoint := network.NewEndpoint(publicAddress, publicPort)

	s.Logger().Debug("nat info",
		"privateAddress", privateAddress,
		"privatePort", privatePort,
		"publicAddress", publicAddress,
		"publicPort", publicPort,
		"unk", unk,
		"connectionType", connectionType,
		"privateEndpoint", privateEndpoint.GetString(),
		"publicEndpoint", publicEndpoint.GetString(),
	)

	if connectionType == 6 {
		connectionType = 4
		s.Logger().Debug("connection type set to 4")
	}

	nat.Set(nat.NatInfo{
		UserId:         s.UserId(),
		PublicAddress:  s.Endpoint().GetAddress(),
		PublicPort:     privatePort,
		PrivateAddress: privateAddress,
		PrivatePort:    privatePort,
		NatUnk:         unk,
		ConnectionType: connectionType,
	})

	return nil
}

func (h *NatInfoReqHandler) AllowAnonymous() bool { return false }

func (h *NatInfoReqHandler) GetHandlerName() string { return "NatInfoReqHandler" }
