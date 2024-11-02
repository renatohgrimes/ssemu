package handlers

import (
	"context"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type ItemBuyReqHandler struct{}

const ItemBuyReq network.PacketId = 0x3B

func (h *ItemBuyReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)
	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()
	ackWriteServerResult(span, &ack, FailedToRequestTask)
	s.Send(&ack)
	return nil
}

func (h *ItemBuyReqHandler) AllowAnonymous() bool { return false }

func (h *ItemBuyReqHandler) GetHandlerName() string { return "ItemBuyReqHandler" }
