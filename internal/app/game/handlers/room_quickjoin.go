package handlers

import (
	"context"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type RoomQuickJoinReqHandler struct{}

const RoomQuickJoinReq network.PacketId = 0x10

func (h *RoomQuickJoinReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)
	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()
	ackWriteServerResult(span, &ack, FailedToRequestTask)
	s.Send(&ack)
	return nil
}

func (h *RoomQuickJoinReqHandler) AllowAnonymous() bool { return false }

func (h *RoomQuickJoinReqHandler) GetHandlerName() string { return "RoomQuickJoinReqHandler" }
