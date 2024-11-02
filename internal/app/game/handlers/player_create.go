package handlers

import (
	"context"
	"database/sql"
	"log/slog"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

type CreatePlayerReqHandler struct{}

const CreatePlayerReq network.PacketId = 0x13

func (h *CreatePlayerReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
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
	if !available {
		ackWriteServerResult(span, &ack, NicknameAlreadyUsed)
		s.Send(&ack)
		return nil
	}

	player := domain.NewPlayer(s.UserId(), nickname)

	tx, err := database.GetConn().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = dbCreatePlayer(ctx, tx, player); err != nil {
		return err
	}

	if err = dbGiveLicense(ctx, tx, player, domain.None); err != nil {
		return err
	}

	if err = dbUpdateUserNickname(ctx, tx, s.UserId(), nickname); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	ackWriteServerResult(span, &ack, PlayerCreateSuccess)
	s.Send(&ack)

	s.Logger().LogAttrs(ctx, slog.LevelDebug, "player created", slog.String("nickname", string(nickname)))

	if err = sendPlayerData(ctx, s); err != nil {
		return err
	}

	return nil
}

func dbCreatePlayer(ctx context.Context, tx *sql.Tx, p domain.Player) error {
	ctx, span := database.GetTracer().Start(ctx, "dbCreatePlayer")
	defer span.End()
	const query string = `
		INSERT INTO players (user_id, nickname, created_utc, tutorial_status)
		VALUES (?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, query,
		p.UserId, p.Nickname, p.CreatedUtc, p.TutorialStatus); err != nil {
		return err
	}
	return nil
}

func dbGiveLicense(ctx context.Context, tx *sql.Tx, player domain.Player, license domain.LicenseId) error {
	ctx, span := database.GetTracer().Start(ctx, "dbGiveLicense")
	defer span.End()
	const query string = "INSERT INTO player_licenses (user_id, license_id) VALUES (?, ?)"
	if _, err := tx.ExecContext(ctx, query, player.UserId, license); err != nil {
		return err
	}
	return nil
}

func dbUpdateUserNickname(ctx context.Context, tx *sql.Tx, u domain.UserId, n domain.Nickname) error {
	ctx, span := database.GetTracer().Start(ctx, "dbUpdateUserNickname")
	defer span.End()
	const query string = `
		UPDATE players
		SET nickname = ?
		WHERE user_id = ?`
	if _, err := tx.ExecContext(ctx, query, n, u); err != nil {
		return err
	}
	return nil
}

func (h *CreatePlayerReqHandler) AllowAnonymous() bool { return false }

func (h *CreatePlayerReqHandler) GetHandlerName() string { return "CreatePlayerReqHandler" }
