package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type LoginReqHandler struct {
	AuthServer network.Server
}

const LoginReq network.PacketId = 0x0A
const LoginAck network.PacketId = 0x0B

type LoginResult byte

const (
	Ok LoginResult = iota
	AccountError
	AccountBlocked
	LoginFailure
)

func (h *LoginReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	usernameStr := p.ReadStringSlice(13)
	passwordStr := p.ReadStringSlice(13)

	isCreateAccount := isCreatingAccount(usernameStr)

	span.SetAttributes(
		attribute.Key("usernameStr").String(usernameStr),
		attribute.Key("isCreateAccount").Bool(isCreateAccount),
	)

	var username domain.Username
	var password domain.Password
	var err error

	ack := network.NewPacket(LoginAck)
	defer ack.Free()

	if isCreateAccount {
		username, err = domain.NewUsername(usernameStr[:len(usernameStr)-2])
	} else {
		username, err = domain.NewUsername(usernameStr)
	}
	if err != nil {
		ackWriteLoginResult(span, &ack, 0, AccountError)
		s.Send(&ack)
		return nil
	}

	if password, err = domain.NewPassword(passwordStr); err != nil {
		ackWriteLoginResult(span, &ack, 0, AccountError)
		s.Send(&ack)
		return nil
	}

	if isCreateAccount {
		user := domain.NewUser(username, password)

		exists, err := dbUserExists(ctx, user.Username)
		if err != nil {
			return err
		}
		if exists {
			ackWriteLoginResult(span, &ack, 0, AccountError)
			s.Send(&ack)
			return nil
		}

		if err = dbCreateUser(ctx, user); err != nil {
			return err
		}

		ackWriteLoginResult(span, &ack, 0, LoginFailure)
		s.Send(&ack)

		s.Logger().LogAttrs(ctx, slog.LevelInfo, "user created",
			slog.Int("userId", int(user.Id)),
			slog.String("username", string(user.Username)),
		)

		return nil
	}

	var user domain.User

	if err := dbGetUserByUsername(ctx, &user, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ackWriteLoginResult(span, &ack, 0, AccountError)
			s.Send(&ack)
			return nil
		}
		return err
	}

	if !user.MatchPassword(passwordStr) {
		ackWriteLoginResult(span, &ack, 0, AccountError)
		s.Send(&ack)
		return nil
	}

	if user.IsBanned() {
		ackWriteLoginResult(span, &ack, 0, AccountBlocked)
		s.Send(&ack)
		return nil
	}

	if concurrentSession, exists := h.AuthServer.GetSessionByUserId(user.Id); exists {
		span.SetAttributes(attribute.Key("concurrentSessionId").Int(int(concurrentSession.Id())))
		concurrentSession.Disconnect("another session connected")
	}

	ackWriteLoginResult(span, &ack, s.Id(), Ok)
	s.Send(&ack)

	s.SetUserId(user.Id)
	span.SetAttributes(attribute.Key("setUserId").Int(int(user.Id)))

	if err := dbUpdateUserLastLogin(ctx, user.Id); err != nil {
		return err
	}

	s.Logger().Info("user logged in")

	return nil
}

func ackWriteLoginResult(span trace.Span, ack *network.Packet, s network.SessionId, r LoginResult) {
	span.SetAttributes(attribute.Key("result").Int(int(r)))
	ack.WriteU32(uint32(s))
	ack.Skip(12)
	ack.WriteU8(byte(r))
}

func isCreatingAccount(username string) bool {
	return strings.HasSuffix(username, "00")
}

func dbUserExists(ctx context.Context, u domain.Username) (bool, error) {
	ctx, span := database.GetTracer().Start(ctx, "dbUserExists")
	defer span.End()
	const query string = "SELECT EXISTS(SELECT 1 FROM users u WHERE u.username = ?)"
	db := database.GetConn()
	r := db.QueryRowContext(ctx, query, u)
	var count int
	if err := r.Scan(&count); err != nil {
		return false, err
	}
	exists := count == 1
	return exists, nil
}

func dbCreateUser(ctx context.Context, u domain.User) error {
	ctx, span := database.GetTracer().Start(ctx, "dbCreateUser")
	defer span.End()
	const query string = `
		INSERT INTO users (id, username, password, created_utc, banned_utc, is_admin, last_login_utc)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	db := database.GetConn()
	if _, err := db.ExecContext(ctx, query,
		u.Id, u.Username, u.Password, u.CreatedUtc, u.BannedUtc, u.IsAdmin, u.LastLoginUtc); err != nil {
		return err
	}
	return nil
}

func dbGetUserByUsername(ctx context.Context, u *domain.User, username domain.Username) error {
	ctx, span := database.GetTracer().Start(ctx, "dbGetByUsername")
	defer span.End()
	const query string = `
		SELECT id, username, password, banned_utc 
		FROM users u
		WHERE u.username = ?
		LIMIT 1`
	db := database.GetConn()
	r := db.QueryRowContext(ctx, query, username)
	if err := r.Scan(&u.Id, &u.Username, &u.Password, &u.BannedUtc); err != nil {
		return err
	}
	return nil
}

func dbUpdateUserLastLogin(ctx context.Context, id domain.UserId) error {
	ctx, span := database.GetTracer().Start(ctx, "dbUpdateLastLogin")
	defer span.End()
	const query string = `
		UPDATE users
		SET last_login_utc = ?
		WHERE id = ?`
	db := database.GetConn()
	if _, err := db.ExecContext(ctx, query, time.Now().UTC(), id); err != nil {
		return err
	}
	return nil
}

func (h *LoginReqHandler) AllowAnonymous() bool { return true }

func (h *LoginReqHandler) GetHandlerName() string { return "LoginReqHandler" }
