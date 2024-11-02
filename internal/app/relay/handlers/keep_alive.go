package handlers

import (
	"context"
	"ssemu/internal/network"
)

type KeepAliveReqHandler struct{}

const KeepAliveReq network.PacketId = 0x01

func (h *KeepAliveReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *KeepAliveReqHandler) AllowAnonymous() bool { return false }

func (h *KeepAliveReqHandler) GetHandlerName() string { return "KeepAliveReqHandler" }
