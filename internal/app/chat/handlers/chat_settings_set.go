package handlers

import (
	"context"
	"ssemu/internal/network"
)

type SetChatSettingsReqHandler struct{}

const SetChatSettingsReq network.PacketId = 0x30

func (h *SetChatSettingsReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *SetChatSettingsReqHandler) AllowAnonymous() bool { return false }

func (h *SetChatSettingsReqHandler) GetHandlerName() string { return "SetChatSettingsReqHandler" }
