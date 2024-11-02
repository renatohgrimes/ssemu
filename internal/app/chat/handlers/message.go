package handlers

import (
	"context"
	"ssemu/internal/network"
)

type MessageReqHandler struct{}

const MessageReq network.PacketId = 0x0A
const MessageAck network.PacketId = 0x0B

func (h *MessageReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *MessageReqHandler) AllowAnonymous() bool { return false }

func (h *MessageReqHandler) GetHandlerName() string { return "MessageReqHandler" }
