package handlers

import (
	"context"
	"ssemu/internal/network"
)

type SetChatStateReqHandler struct{}

const SetChatStateReq network.PacketId = 0x33

func (h *SetChatStateReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *SetChatStateReqHandler) AllowAnonymous() bool { return false }

func (h *SetChatStateReqHandler) GetHandlerName() string { return "SetChatStateReqHandler" }
