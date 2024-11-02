package handlers

import (
	"context"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"ssemu/internal/utils"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type LoginReqHandler struct {
	RelayServer network.Server
	AuthServer  network.Server
	GameServer  network.Server
}

const LoginReq network.PacketId = 0x03
const LoginAck network.PacketId = 0x02

func (h *LoginReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	nicknameStr := p.ReadStringSlice(16)

	span.SetAttributes(attribute.String("nicknameStr", nicknameStr))

	nickname, err := domain.NewNickname(nicknameStr)
	if err != nil {
		return err
	}

	userId, err := dbGetUserIdByNickname(ctx, nickname)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	if _, exists := h.GameServer.GetSessionByUserId(userId); !exists {
		if _, exists := h.AuthServer.GetSessionByUserId(userId); !exists {
			return utils.NewValidationError("session not found")
		}
	}

	if _, exists := h.RelayServer.GetSessionByUserId(userId); exists {
		return utils.NewValidationError("duplicate relay session")
	}

	s.SetUserId(userId)
	span.SetAttributes(attribute.Int("setUserId", int(userId)))

	if gameSession, exists := h.GameServer.GetSessionByUserId(userId); exists {
		gameSession.SetOnDisconnectCallback(func() {
			time.Sleep(time.Second)
			s.Disconnect("game session disconnected")
		})
	}

	ack := network.NewPacket(LoginAck)
	defer ack.Free()

	ack.WriteU32(0) // no errors

	s.Send(&ack)

	s.Logger().Debug("user logged in")

	return nil
}

func dbGetUserIdByNickname(ctx context.Context, nickname domain.Nickname) (domain.UserId, error) {
	ctx, span := database.GetTracer().Start(ctx, "dbGetUserIdByNickname")
	defer span.End()
	const query string = "SELECT user_id FROM players p WHERE p.nickname = ? LIMIT 1"
	db := database.GetConn()
	r := db.QueryRowContext(ctx, query, nickname)
	var id domain.UserId
	if err := r.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (h *LoginReqHandler) AllowAnonymous() bool { return true }

func (h *LoginReqHandler) GetHandlerName() string { return "LoginReqHandler" }
