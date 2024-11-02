package handlers

import (
	"context"
	"ssemu/internal/network"
)

type WhisperReqHandler struct{}

const WhisperReq network.PacketId = 0x10
const WhisperAck network.PacketId = 0x11

func (h *WhisperReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *WhisperReqHandler) AllowAnonymous() bool { return false }

func (h *WhisperReqHandler) GetHandlerName() string { return "WhisperReqHandler" }
