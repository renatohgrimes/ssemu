package handlers

import (
	"context"
	"ssemu/internal/network"
)

type TunnelJoinReqHandler struct{}

const TunnelJoinReq network.PacketId = 0x1C
const TunnelJoinAck network.PacketId = 0x1D

func (h *TunnelJoinReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	// TODO do we actually need this?
	slot := p.ReadU8()
	// TODO validate if user is inside a room
	ack := network.NewPacket(TunnelJoinAck)
	defer ack.Free()
	ack.WriteU8(slot)
	s.Send(&ack)
	return nil
}

func (h *TunnelJoinReqHandler) AllowAnonymous() bool { return false }

func (h *TunnelJoinReqHandler) GetHandlerName() string { return "TunnelJoinReqHandler" }
