package handlers

import (
	"context"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type RefundItemReqHandler struct{}

const RefundItemReq network.PacketId = 0x41

func (h *RefundItemReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)
	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()
	ackWriteServerResult(span, &ack, FailedToRequestTask)
	s.Send(&ack)
	return nil
}

func (h *RefundItemReqHandler) AllowAnonymous() bool { return false }

func (h *RefundItemReqHandler) GetHandlerName() string { return "RefundItemReqHandler" }
