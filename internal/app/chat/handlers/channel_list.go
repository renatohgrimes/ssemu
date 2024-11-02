package handlers

import (
	"context"
	"ssemu/internal/network"
)

type ChannelListReqHandler struct{}

const ChannelListReq network.PacketId = 0x0C
const ChannelListAck network.PacketId = 0x0D

func (h *ChannelListReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *ChannelListReqHandler) AllowAnonymous() bool { return false }

func (h *ChannelListReqHandler) GetHandlerName() string { return "ChannelListReqHandler" }
