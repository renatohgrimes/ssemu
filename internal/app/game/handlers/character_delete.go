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

type CharacterDeleteReqHandler struct{}

const CharacterDeleteReq network.PacketId = 0x0B
const CharacterDeleteAck network.PacketId = 0x0C

func (h *CharacterDeleteReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	slot := p.ReadU8()

	span.SetAttributes(attribute.Int("slot", int(slot)))

	charCount, activeSlot, err := dbGetCharacterSlotData(ctx, s.UserId())
	if err != nil {
		return err
	}

	if charCount == 1 {
		return utils.NewValidationError("user must have at least one character")
	}

	if slot == activeSlot {
		return utils.NewValidationError("trying to delete active character")
	}

	exists, err := dbIsCharacterSlotExist(ctx, s.UserId(), slot)
	if err != nil {
		return err
	}
	if !exists {
		return utils.NewValidationError("trying to delete an invalid character slot")
	}

	if err := dbDeleteCharacter(ctx, s.UserId(), slot); err != nil {
		return err
	}

	ack := network.NewPacket(CharacterDeleteAck)
	defer ack.Free()
	ack.WriteU8(slot)

	s.Send(&ack)

	return nil
}

func dbDeleteCharacter(ctx context.Context, uid domain.UserId, slot byte) error {
	ctx, span := database.GetTracer().Start(ctx, "dbDeleteCharacter")
	defer span.End()
	const cmd string = "DELETE FROM player_characters WHERE user_id = ? AND slot = ?"
	if _, err := database.GetConn().ExecContext(ctx, cmd, uid, slot); err != nil {
		return err
	}
	return nil
}

func (h *CharacterDeleteReqHandler) AllowAnonymous() bool { return false }

func (h *CharacterDeleteReqHandler) GetHandlerName() string { return "CharacterDeleteReqHandler" }
