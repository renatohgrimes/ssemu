package handlers

import (
	"context"
	"ssemu/internal/network"
)

type GetChatSettingsReqHandler struct{}

const GetChatSettingsReq network.PacketId = 0x31
const GetChatSettingsAck network.PacketId = 0x32

func (h *GetChatSettingsReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *GetChatSettingsReqHandler) AllowAnonymous() bool { return false }

func (h *GetChatSettingsReqHandler) GetHandlerName() string { return "GetChatSettingsReqHandler" }
