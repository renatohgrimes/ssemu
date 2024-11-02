package network

import (
	"context"
	"log/slog"
	"ssemu/internal/domain"
)

var logLevelPacket slog.Level = slog.LevelDebug - 4

type ServerType byte

const (
	Auth ServerType = iota
	Game
	Chat
	Nat
	Relay
)

type ServerSettings struct {
	Name        string
	Port        uint16
	Type        ServerType
	MaxCapacity int
	Version     string
}

type packetHandler interface {
	HandlePacket(context.Context, Session, Packet) error
	AllowAnonymous() bool
	GetHandlerName() string
}

type Server interface {
	ListenAndServe()
	Shutdown()
	RegisterHandler(PacketId, packetHandler) error
	GetSessionCount() int
	GetPublicEndpoint() Endpoint
	GetSessionById(SessionId) (Session, bool)
	GetSessionByUserId(domain.UserId) (Session, bool)
}
