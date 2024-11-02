package handlers

import (
	"context"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type RepairItemReqHandler struct{}

const RepairItemReq network.PacketId = 0x3D

func (h *RepairItemReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)
	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()
	ackWriteServerResult(span, &ack, FailedToRequestTask)
	s.Send(&ack)
	return nil
}

func (h *RepairItemReqHandler) AllowAnonymous() bool { return false }

func (h *RepairItemReqHandler) GetHandlerName() string { return "RepairItemReqHandler" }
