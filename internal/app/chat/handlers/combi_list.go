package handlers

import (
	"context"
	"ssemu/internal/network"
)

type CombiListReqHandler struct{}

const CombiListReq network.PacketId = 0x3A
const CombiListAck network.PacketId = 0x3B

func (h *CombiListReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *CombiListReqHandler) AllowAnonymous() bool { return false }

func (h *CombiListReqHandler) GetHandlerName() string { return "CombiListReqHandler" }
