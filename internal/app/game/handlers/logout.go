package handlers

import (
	"context"
	"ssemu/internal/network"
	"time"
)

type LogoutReqHandler struct{}

const LogoutReq network.PacketId = 0x5C
const LogoutAck network.PacketId = 0x5D

func (h *LogoutReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	ack := network.NewPacket(LogoutAck)
	defer ack.Free()
	s.Send(&ack)
	time.Sleep(time.Second)
	s.Disconnect("user sent a logout request")
	return nil
}

func (h *LogoutReqHandler) AllowAnonymous() bool { return false }

func (h *LogoutReqHandler) GetHandlerName() string { return "LogoutReqHandler" }
