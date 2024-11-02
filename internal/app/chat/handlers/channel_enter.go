package handlers

import (
	"context"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/app/game/nat"
	"ssemu/internal/network"
	"ssemu/internal/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ChannelEnterReqHandler struct{}

const ChannelEnterReq network.PacketId = 0x05
const ChannelEnterAck network.PacketId = 0x06

func (h *ChannelEnterReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	channelStr := p.ReadStringSlice(20)

	span.SetAttributes(attribute.String("channelStr", channelStr))

	channel, err := channels.GetChannelByName(channelStr)
	if err != nil {
		return utils.NewValidationError("channel not found")
	}

	ack := network.NewPacket(ChannelEnterAck)
	defer ack.Free()

	ack.WriteU32(channel.GetId())

	s.Send(&ack)

	if _, err := nat.Get(s.UserId()); err != nil {
		msg := "Unknown NAT info from your connection. You cannot join rooms."
		ack2 := network.NewPacket(MessageAck)
		defer ack2.Free()
		ack2.WriteU64(uint64(s.UserId()))
		ack2.WriteU32(channel.GetId())
		ack2.WriteU16(uint16(len(msg)))
		ack2.WriteStringSlice(msg, len(msg))
		s.Send(&ack2)
	}

	return nil
}

func (h *ChannelEnterReqHandler) AllowAnonymous() bool { return false }

func (h *ChannelEnterReqHandler) GetHandlerName() string { return "ChannelEnterReqHandler" }
