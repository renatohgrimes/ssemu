package network

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"ssemu/internal/domain"
	"time"
)

type SessionId uint32

type Session interface {
	Id() SessionId
	Disconnect(string)
	SetUserId(domain.UserId)
	IsAuthenticated() bool
	Send(*Packet)
	Endpoint() Endpoint
	UserId() domain.UserId
	Logger() *slog.Logger
	SetOnDisconnectCallback(func())
}

type session struct {
	id                  SessionId
	userId              domain.UserId
	connection          net.Conn
	remoteEndpoint      Endpoint
	listening           bool
	context             context.Context
	cancelContext       context.CancelFunc
	logger              *slog.Logger
	disconnectCallbacks []func()
}

func newSession(ctx context.Context, id SessionId, conn net.Conn, logger *slog.Logger) *session {
	ctx, cancelFunc := context.WithCancel(ctx)
	remoteEndpoint := NewEndpointFromAddress(conn.RemoteAddr())
	return &session{
		id:             id,
		userId:         0,
		connection:     conn,
		remoteEndpoint: remoteEndpoint,
		listening:      true,
		context:        ctx,
		cancelContext:  cancelFunc,
		logger: logger.With(
			slog.Int("session", int(id)),
			slog.String("remote", remoteEndpoint.GetString()),
		),
		disconnectCallbacks: make([]func(), 0, 3),
	}
}

func (s *session) Send(p *Packet) {
	if !s.listening {
		s.logger.Debug("session not listening")
		return
	}
	buffer := p.Data()
	deadline := time.Now().Add(5 * time.Second)
	_ = s.connection.SetWriteDeadline(deadline)
	if _, err := s.connection.Write(buffer); err != nil {
		s.Disconnect("failed to write buffer")
		s.logger.LogAttrs(s.context, slog.LevelWarn, "failed to write buffer",
			slog.String("packet", p.id.HexString()),
			slog.Int("length", len(buffer)),
			slog.String("err", err.Error()),
		)
		return
	}
	s.logger.LogAttrs(s.context, logLevelPacket, "packet sent",
		slog.String("packet", p.id.HexString()),
		slog.Int("length", len(buffer)),
	)
}

func (s *session) Receive(buffer []byte) (int, error) {
	if !s.listening {
		return 0, errors.New("session not listening")
	}
	_ = s.connection.SetReadDeadline(time.Now().Add(30 * time.Second))
	return s.connection.Read(buffer)
}

func (s *session) Disconnect(reason string) {
	if !s.listening {
		return
	}
	s.listening = false
	s.cancelContext()
	_ = s.connection.Close()
	for _, callback := range s.disconnectCallbacks {
		callback()
	}
	s.logger.LogAttrs(s.context, slog.LevelDebug, "session disconnected",
		slog.String("reason", reason),
	)
}

func (s *session) SetUserId(id domain.UserId) {
	s.userId = id
	s.logger = s.logger.With(slog.Int("user", int(s.userId)))
}

func (s *session) SetOnDisconnectCallback(callback func()) {
	s.disconnectCallbacks = append(s.disconnectCallbacks, callback)
}

func (s session) IsAuthenticated() bool { return s.userId != 0 }

func (s session) Id() SessionId { return s.id }

func (s session) Endpoint() Endpoint { return s.remoteEndpoint }

func (s session) Context() context.Context { return s.context }

func (s session) UserId() domain.UserId { return s.userId }

func (s session) Logger() *slog.Logger { return s.logger }
