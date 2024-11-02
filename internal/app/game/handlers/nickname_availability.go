package handlers

import (
	"context"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type GetNicknameAvailabilityReqHandler struct{}

const GetNicknameAvailabilityReq network.PacketId = 0x14

func (h *GetNicknameAvailabilityReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	nicknameStr := p.ReadStringSlice(16)

	ack := network.NewPacket(ServerResultAck)
	defer ack.Free()

	nickname, err := domain.NewNickname(nicknameStr)
	if err != nil {
		return err
	}

	available, err := dbIsNicknameAvailable(ctx, nickname)
	if err != nil {
		return err
	}
	if available {
		ackWriteServerResult(span, &ack, NicknameAvailable)
	} else {
		ackWriteServerResult(span, &ack, NicknameAlreadyUsed)
	}

	s.Send(&ack)

	return nil
}

func dbIsNicknameAvailable(ctx context.Context, nickname domain.Nickname) (bool, error) {
	ctx, span := database.GetTracer().Start(ctx, "dbIsNicknameAvailable")
	defer span.End()
	const query string = "SELECT EXISTS(SELECT 1 FROM players WHERE nickname = ?)"
	db := database.GetConn()
	r := db.QueryRowContext(ctx, query, nickname)
	var exists bool
	if err := r.Scan(&exists); err != nil {
		return false, err
	}
	return !exists, nil
}

func (h *GetNicknameAvailabilityReqHandler) AllowAnonymous() bool { return false }

func (h *GetNicknameAvailabilityReqHandler) GetHandlerName() string {
	return "GetNicknameAvailabilityReqHandler"
}
