package handlers

import (
	"context"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type TunnelLeaveReqHandler struct{}

const TunnelLeaveReq network.PacketId = 0x0A

func (h *TunnelLeaveReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)
	channel, err := channels.GetUserChannel(s.UserId())
	if err != nil {
		return err
	}
	// TODO room must be its own entity, not being controlled by the channel
	channel.LeaveRoom(s.UserId())
	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()
	ackWriteServerResult(span, &ack, TunnelLeaveSuccess)
	s.Send(&ack)
	return nil
}

func (h *TunnelLeaveReqHandler) AllowAnonymous() bool { return false }

func (h *TunnelLeaveReqHandler) GetHandlerName() string { return "TunnelLeaveReqHandler" }
