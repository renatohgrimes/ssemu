package handlers

import (
	"context"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"ssemu/internal/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CharacterSelectReqHandler struct{}

const CharacterSelectReq network.PacketId = 0x0D
const CharacterSelectAck network.PacketId = 0x0E

func (h *CharacterSelectReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	slot := p.ReadU8()

	span.SetAttributes(attribute.Int("slot", int(slot)))

	ack := network.NewPacket(CharacterSelectAck)
	defer ack.Free()

	exists, err := dbIsCharacterSlotExist(ctx, s.UserId(), slot)
	if err != nil {
		return err
	} else if !exists {
		return utils.NewValidationError("character not found")
	}

	if err := dbChangeCharacter(ctx, s.UserId(), slot); err != nil {
		return err
	}

	span.SetAttributes(attribute.Bool("exists", true))
	ack.WriteU8(slot)
	s.Send(&ack)

	return nil
}

func dbIsCharacterSlotExist(ctx context.Context, uid domain.UserId, slot byte) (bool, error) {
	ctx, span := database.GetTracer().Start(ctx, "dbIsCharacterSlotExist")
	defer span.End()
	const query string = "SELECT EXISTS(SELECT 1 FROM player_characters WHERE user_id = ? AND slot = ?)"
	db := database.GetConn()
	r := db.QueryRowContext(ctx, query, uid, slot)
	var exists bool
	if err := r.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func dbChangeCharacter(ctx context.Context, uid domain.UserId, slot byte) error {
	ctx, span := database.GetTracer().Start(ctx, "dbChangeCharacter")
	defer span.End()
	const cmd1 string = "UPDATE player_characters SET is_active = true WHERE user_id = ? AND slot = ?"
	const cmd2 string = "UPDATE player_characters SET is_active = false WHERE user_id = ? AND slot <> ?"
	tx, err := database.GetConn().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, cmd1, uid, slot); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, cmd2, uid, slot); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (h *CharacterSelectReqHandler) AllowAnonymous() bool { return false }

func (h *CharacterSelectReqHandler) GetHandlerName() string { return "CharacterSelectReqHandler" }
