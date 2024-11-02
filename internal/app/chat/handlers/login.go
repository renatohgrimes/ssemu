package handlers

import (
	"context"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"ssemu/internal/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type LoginReqHandler struct {
	GameServer network.Server
	ChatServer network.Server
}

const LoginReq network.PacketId = 0x03
const LoginAck network.PacketId = 0x02

func (h *LoginReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	userId := domain.UserId(p.ReadU32())

	gameSession, exists := h.GameServer.GetSessionByUserId(userId)
	if !exists {
		return utils.NewValidationError("game session does not exists")
	}

	if _, exists := h.ChatServer.GetSessionByUserId(userId); exists {
		return utils.NewValidationError("chat user already exists")
	}

	s.SetUserId(userId)
	span.SetAttributes(attribute.Key("setUserId").Int(int(userId)))

	gameSession.SetOnDisconnectCallback(func() {
		s.Disconnect("game session disconnected")
	})

	ack := network.NewPacket(LoginAck)
	defer ack.Free()
	ack.WriteU32(0) // no errors

	s.Send(&ack)

	s.Logger().Debug("user logged in")

	return nil
}

func (h *LoginReqHandler) AllowAnonymous() bool { return true }

func (h *LoginReqHandler) GetHandlerName() string { return "LoginReqHandler" }
