package handlers

import (
	"context"
	"errors"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/domain"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ChannelEnterReqHandler struct{}

const ChannelEnterReq network.PacketId = 0x2B
const RefreshMoneyAck network.PacketId = 0x43

func (h *ChannelEnterReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	chanId := p.ReadU32()

	span.SetAttributes(attribute.Int("chanId", int(chanId)))

	channel, err := channels.GetChannelById(chanId)
	if err != nil {
		return err
	}

	err = channel.Join(s)
	if err != nil {
		if errors.Is(err, channels.ErrChannelCapacityExceeded) {
			ack := network.NewPacket(ServerResultAck)
			defer ack.Free()
			ackWriteServerResult(span, &ack, ChannelCapacityExceeded)
			s.Send(&ack)
			return nil
		}
		return err
	}

	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()

	ackWriteServerResult(span, &ack, ChannelEnter)

	s.Send(&ack)

	ack2 := network.NewPacket(RefreshMoneyAck)
	defer ack2.Free()

	ack2.WriteU32(domain.Mocks.Pen)
	ack2.WriteU32(domain.Mocks.Cash)

	s.Send(&ack2)

	return nil
}

func (h *ChannelEnterReqHandler) AllowAnonymous() bool { return false }

func (h *ChannelEnterReqHandler) GetHandlerName() string { return "ChannelEnterReqHandler" }
