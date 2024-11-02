package handlers

import (
	"context"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type TunnelJoinReqHandler struct{}

const TunnelJoinReq network.PacketId = 0x04

func (h *TunnelJoinReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)
	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()
	ackWriteServerResult(span, &ack, TunnelJoinSuccess)
	s.Send(&ack)
	return nil
}

func (h *TunnelJoinReqHandler) AllowAnonymous() bool { return false }

func (h *TunnelJoinReqHandler) GetHandlerName() string { return "TunnelJoinReqHandler" }
