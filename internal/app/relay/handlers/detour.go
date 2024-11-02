package handlers

import (
	"context"
	"ssemu/internal/network"
)

type DetourReqHandler struct {
	RelayServer network.Server
}

const DetourReq network.PacketId = 0x07
const DetourAck network.PacketId = 0x09

func (h *DetourReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *DetourReqHandler) AllowAnonymous() bool { return false }

func (h *DetourReqHandler) GetHandlerName() string { return "DetourReqHandler" }
