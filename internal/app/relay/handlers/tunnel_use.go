package handlers

import (
	"context"
	"ssemu/internal/network"
)

type TunnelUseReqHandler struct{}

const TunnelUseReq network.PacketId = 0x05
const TunnelUseAck network.PacketId = 0x06

func (h *TunnelUseReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	slotId := p.ReadU8()
	// TODO validate slot ?
	ack := network.NewPacket(TunnelUseAck)
	defer ack.Free()
	ack.WriteU8(slotId)
	s.Send(&ack)
	return nil
}

func (h *TunnelUseReqHandler) AllowAnonymous() bool { return false }

func (h *TunnelUseReqHandler) GetHandlerName() string { return "TunnelUseReqHandler" }
