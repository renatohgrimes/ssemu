package handlers

import (
	"context"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type ChannelLeaveReqHandler struct{}

const ChannelLeaveReq network.PacketId = 0x2C

func (h *ChannelLeaveReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)
	channels.RemoveUser(s.UserId())
	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()
	ackWriteServerResult(span, &ack, ChannelLeave)
	s.Send(&ack)
	return nil
}

func (h *ChannelLeaveReqHandler) AllowAnonymous() bool { return false }

func (h *ChannelLeaveReqHandler) GetHandlerName() string { return "ChannelLeaveReqHandler" }
