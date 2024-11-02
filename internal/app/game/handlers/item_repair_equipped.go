package handlers

import (
	"context"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type RepairEquippedItemReqHandler struct{}

const RepairEquippedItemReq network.PacketId = 0x3E

func (h *RepairEquippedItemReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)
	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()
	ackWriteServerResult(span, &ack, FailedToRequestTask)
	s.Send(&ack)
	return nil
}

func (h *RepairEquippedItemReqHandler) AllowAnonymous() bool { return false }

func (h *RepairEquippedItemReqHandler) GetHandlerName() string { return "RepairEquippedItemReqHandler" }
