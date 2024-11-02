package handlers

import (
	"context"
	"errors"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type RoomJoinReqHandler struct {
}

const RoomJoinReq network.PacketId = 0x36

func (h *RoomJoinReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	roomId := p.ReadU32()
	password := p.ReadU32()

	channel, err := channels.GetUserChannel(s.UserId())
	if err != nil {
		return err
	}

	if err := channel.JoinRoom(ctx, s, roomId, password); err != nil {
		ack := network.NewPacket(ServerResultAck)
		defer ack.Free()
		if errors.Is(err, channels.ErrRoomWrongPassword) || errors.Is(err, channels.ErrRoomCannotEnter) {
			ackWriteServerResult(span, &ack, ImpossibleToEnterRoom)
			s.Send(&ack)
			return nil
		} else if errors.Is(err, channels.ErrRoomCapacityExceed) {
			ackWriteServerResult(span, &ack, RoomCapacityExceed)
			s.Send(&ack)
			return nil
		} else if errors.Is(err, channels.ErrRoomChangingRules) {
			ackWriteServerResult(span, &ack, RoomChangingRules)
			s.Send(&ack)
			return nil
		}
		return err
	}

	return nil
}

func (h *RoomJoinReqHandler) AllowAnonymous() bool { return false }

func (h *RoomJoinReqHandler) GetHandlerName() string { return "RoomJoinReqHandler" }
