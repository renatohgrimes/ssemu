package handlers

import (
	"context"
	"ssemu/internal/network"
)

type ChannelLeaveReqHandler struct{}

const ChannelLeaveReq network.PacketId = 0x04

func (h *ChannelLeaveReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *ChannelLeaveReqHandler) AllowAnonymous() bool { return false }

func (h *ChannelLeaveReqHandler) GetHandlerName() string { return "ChannelLeaveReqHandler" }
