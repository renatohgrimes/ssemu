package network

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"ssemu/internal/domain"
	"ssemu/internal/telemetry"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

type tcpServer struct {
	listener       net.Listener
	listening      bool
	name           string
	sessions       map[SessionId]Session
	handlers       map[PacketId]packetHandler
	publicEndpoint Endpoint
	serverType     ServerType
	context        context.Context
	hostname       string
	tracerProvider *sdktrace.TracerProvider
	rng            *rand.Rand
	logger         *slog.Logger
}

func NewTcpServer(context context.Context, publicIpv4 string, s ServerSettings) (Server, error) {
	if s.Type == Nat {
		return nil, errors.New("unsupported server type protocol")
	}

	logger := slog.Default().With(slog.String("server", s.Name))
	logger.Debug("loading...")

	address := fmt.Sprintf("0.0.0.0:%d", s.Port)
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	tp, err := telemetry.NewTraceProvider(context, telemetry.NewResource(s.Name, s.Version))
	if err != nil {
		return nil, err
	}

	return &tcpServer{
		listener:       listener,
		name:           s.Name,
		listening:      false,
		sessions:       make(map[SessionId]Session),
		handlers:       make(map[PacketId]packetHandler),
		publicEndpoint: NewEndpointFromString(publicIpv4, s.Port),
		serverType:     s.Type,
		context:        context,
		hostname:       hostname,
		tracerProvider: tp,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
		logger:         logger,
	}, nil
}

func (s *tcpServer) ListenAndServe() {
	s.listening = true
	s.logger.LogAttrs(s.context, slog.LevelInfo, "listening...",
		slog.String("endpoint", s.listener.Addr().String()),
		slog.String("public", s.publicEndpoint.GetString()),
		slog.String("protocol", "tcp4"),
	)
	for {
		conn, _ := s.listener.Accept()
		if !s.listening {
			return
		}
		session := s.createSession(conn)
		go s.listenSession(session)
	}
}

func (s *tcpServer) createSession(conn net.Conn) *session {
	sessionId := s.generateSessionId()
	session := newSession(s.context, sessionId, conn, s.logger)
	s.sessions[sessionId] = session
	session.logger.Debug("session created")
	return session
}

func (s *tcpServer) generateSessionId() SessionId {
	for {
		id := SessionId(s.rng.Uint32())
		if _, exists := s.sessions[id]; !exists {
			return id
		}
	}
}

func (s *tcpServer) listenSession(session *session) {
	buffer := make([]byte, packetBufferSize)
	limiter := rate.NewLimiter(10, 1)
	for {
		if !session.listening || !s.listening {
			break
		}
		read, _ := session.Receive(buffer)
		if read < 4 {
			break
		}
		if s.serverType == Game {
			Decrypt(buffer[4:read])
		}
		packet := GetPacket(buffer[:read])
		session.logger.LogAttrs(session.context, logLevelPacket, "packet received",
			slog.String("packet", packet.id.HexString()),
			slog.Int("length", read),
		)
		handler, exists := s.handlers[packet.id]
		if !exists {
			session.Disconnect("handler not found")
			session.logger.LogAttrs(session.context, slog.LevelError, "handler not found",
				slog.String("packet", packet.id.HexString()),
			)
			break
		}
		if !handler.AllowAnonymous() && !session.IsAuthenticated() {
			session.Disconnect("session not authenticated")
			session.logger.LogAttrs(session.context, slog.LevelError, "session not authenticated",
				slog.String("packet", packet.id.HexString()),
			)
			break
		}
		if err := s.safeHandlePacket(session.context, session, handler, packet); err != nil {
			break
		}
		if err := limiter.Wait(session.context); err != nil {
			break
		}
	}
	session.Disconnect("receive loop break")
	delete(s.sessions, session.id)
}

func (s *tcpServer) safeHandlePacket(ctx context.Context, session *session, handler packetHandler, packet Packet) error {
	ctx, span := s.tracerProvider.Tracer("").Start(ctx, handler.GetHandlerName())
	defer span.End()
	span.SetAttributes(
		attribute.Key("user").Int(int(session.UserId())),
		attribute.Key("session").Int(int(session.Id())),
		attribute.Key("remote").String(session.remoteEndpoint.GetString()),
		attribute.Key("packet").Int(int(packet.id)),
		attribute.Key("hostname").String(s.hostname),
		attribute.Key("server").String(s.name),
	)
	defer s.recoverPanic(span, session, packet, handler)
	if err := handler.HandlePacket(ctx, session, packet); err != nil {
		session.Disconnect(err.Error())
		span.SetStatus(codes.Error, "handler error")
		span.RecordError(err)
		session.logger.LogAttrs(ctx, slog.LevelError, "handler error",
			slog.String("handler", handler.GetHandlerName()),
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (s *tcpServer) recoverPanic(span trace.Span, session Session, packet Packet, handler packetHandler) {
	if r := recover(); r != nil {
		err := fmt.Errorf("handler panic: %s", r)
		session.Disconnect(err.Error())
		span.SetStatus(codes.Error, "handler panic")
		span.RecordError(err)
		session.Logger().LogAttrs(s.context, slog.LevelError, "handler panic",
			slog.String("handler", handler.GetHandlerName()),
			slog.String("error", err.Error()),
		)
		packet.DumpPacket(s.name)
	}
}

func (s *tcpServer) RegisterHandler(id PacketId, handler packetHandler) error {
	if _, exists := s.handlers[id]; exists {
		return fmt.Errorf("packet %s handler already registered", id.HexString())
	}
	s.handlers[id] = handler
	return nil
}

func (s *tcpServer) GetSessionById(id SessionId) (Session, bool) {
	session, exists := s.sessions[id]
	return session, exists
}

func (s *tcpServer) GetSessionByUserId(userId domain.UserId) (Session, bool) {
	for _, session := range s.sessions {
		if session.UserId() == userId {
			return session, true
		}
	}
	return nil, false
}

func (s *tcpServer) Shutdown() {
	if !s.listening {
		return
	}

	s.logger.Info("stopping...")

	s.listening = false
	s.listener.Close()
	s.tracerProvider.Shutdown(s.context)

	s.logger.Info("server closed")
}

func (s *tcpServer) GetSessionCount() int { return len(s.sessions) }

func (s *tcpServer) GetPublicEndpoint() Endpoint { return s.publicEndpoint }
