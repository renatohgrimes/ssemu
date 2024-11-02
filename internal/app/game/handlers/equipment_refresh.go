package handlers

import (
	"context"
	"ssemu/internal/network"
)

type RefreshEquipmentReqHandler struct{}

const RefreshEquipmentReq network.PacketId = 0x75
const RefreshEquipmentAck network.PacketId = 0x76

func (h *RefreshEquipmentReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	ack := network.NewPacket(RefreshEquipmentAck)
	defer ack.Free()
	ack.WriteU8(0) // count
	s.Send(&ack)
	return nil
}

func (h *RefreshEquipmentReqHandler) AllowAnonymous() bool { return false }

func (h *RefreshEquipmentReqHandler) GetHandlerName() string { return "RefreshEquipmentReqHandler" }
