package handlers

import (
	"context"
	"ssemu/internal/network"
)

type TimeSyncReqHandler struct {
}

const TimeSyncReq network.PacketId = 0x1E
const TimeSyncAck network.PacketId = 0x1F

func (h *TimeSyncReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *TimeSyncReqHandler) AllowAnonymous() bool { return false }

func (h *TimeSyncReqHandler) GetHandlerName() string { return "TimeSyncReqHandler" }
