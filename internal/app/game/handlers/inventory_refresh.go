package handlers

import (
	"context"
	"ssemu/internal/network"
)

type RefreshInventoryReqHandler struct{}

const RefreshInventoryReq network.PacketId = 0x4F
const RefreshInvalidItemsAck network.PacketId = 0x50

func (h *RefreshInventoryReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	ack := network.NewPacket(RefreshInvalidItemsAck)
	defer ack.Free()
	ack.WriteU8(0) // count
	s.Send(&ack)
	return nil
}

func (h *RefreshInventoryReqHandler) AllowAnonymous() bool { return false }

func (h *RefreshInventoryReqHandler) GetHandlerName() string { return "RefreshInventoryReqHandler" }
