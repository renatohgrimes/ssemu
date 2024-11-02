package handlers

import (
	"context"
	"log/slog"
	"ssemu/internal/app"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"ssemu/internal/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type LoginReqHandler struct {
	AuthServer network.Server
	GameServer network.Server
}

const LoginReq network.PacketId = 0x03
const LoginAck network.PacketId = 0x04

type LoginResult uint32

const (
	Ok                  LoginResult = 0
	PlayerLimitExceeded LoginResult = 1
	AbnormalConnection  LoginResult = 2
	NicknamePending     LoginResult = 4
)

func (h *LoginReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	ack := network.NewPacket(LoginAck)
	defer ack.Free()

	if h.GameServer.GetSessionCount() > app.MaxCapacity {
		ackWriteLoginResult(span, &ack, 0, PlayerLimitExceeded)
		s.Send(&ack)
		return nil
	}

	usernameStr := p.ReadStringSlice(13)
	p.Skip(30)
	authSessionId := network.SessionId(p.ReadU32())

	span.SetAttributes(
		attribute.Key("usernameStr").String(usernameStr),
		attribute.Key("authSessionId").Int(int(authSessionId)))

	if _, exists := h.AuthServer.GetSessionById(authSessionId); !exists {
		return utils.NewValidationError("unknown auth session")
	}

	username, err := domain.NewUsername(usernameStr)
	if err != nil {
		return err
	}

	userId, err := dbGetUserIdByUsername(ctx, username)
	if err != nil {
		return err
	}

	if _, exists := h.GameServer.GetSessionByUserId(userId); exists {
		ackWriteLoginResult(span, &ack, 0, AbnormalConnection)
		s.Send(&ack)
		return nil
	}

	s.SetUserId(userId)
	span.SetAttributes(attribute.Key("setUserId").Int(int(userId)))

	var pending bool
	if pending, err = dbIsNicknamePending(ctx, userId); err != nil {
		return err
	}
	if pending {
		ackWriteLoginResult(span, &ack, userId, NicknamePending)
		s.Send(&ack)
		return nil
	}

	s.SetOnDisconnectCallback(func() { channels.RemoveUser(s.UserId()) })

	ackWriteLoginResult(span, &ack, userId, Ok)
	s.Send(&ack)

	if err = sendPlayerData(ctx, s); err != nil {
		return err
	}

	s.Logger().LogAttrs(ctx, slog.LevelDebug, "user logged in", slog.String("username", string(username)))

	return nil
}

func ackWriteLoginResult(span trace.Span, ack *network.Packet, id domain.UserId, res LoginResult) {
	span.SetAttributes(attribute.Key("result").Int(int(res)))
	ack.WriteU32(uint32(id))
	ack.Skip(4)
	ack.WriteU32(uint32(res))
}

func dbGetUserIdByUsername(ctx context.Context, username domain.Username) (domain.UserId, error) {
	ctx, span := database.GetTracer().Start(ctx, "dbGetUserIdByUsername")
	defer span.End()
	const query string = "SELECT id FROM users u WHERE u.username = ? LIMIT 1"
	db := database.GetConn()
	r := db.QueryRowContext(ctx, query, username)
	var userId uint32
	if err := r.Scan(&userId); err != nil {
		return 0, err
	}
	return domain.UserId(userId), nil
}

func dbIsNicknamePending(ctx context.Context, id domain.UserId) (bool, error) {
	ctx, span := database.GetTracer().Start(ctx, "dbIsNicknamePending")
	defer span.End()
	const query string = "SELECT EXISTS(SELECT 1 FROM players p WHERE p.user_id = ? AND p.nickname IS NOT NULL)"
	db := database.GetConn()
	r := db.QueryRowContext(ctx, query, id)
	var nicknameSet bool
	if err := r.Scan(&nicknameSet); err != nil {
		return true, err
	}
	return !nicknameSet, nil
}

func (h *LoginReqHandler) AllowAnonymous() bool { return true }

func (h *LoginReqHandler) GetHandlerName() string { return "LoginReqHandler" }
