package handlers

import (
	"context"
	"ssemu/internal/network"
)

type AdminCommandReqHandler struct{}

const AdminCommandReq network.PacketId = 0x44

func (h *AdminCommandReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *AdminCommandReqHandler) AllowAnonymous() bool { return false }

func (h *AdminCommandReqHandler) GetHandlerName() string { return "AdminCommandReqHandler" }
