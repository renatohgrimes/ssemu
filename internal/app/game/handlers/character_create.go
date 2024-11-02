package handlers

import (
	"context"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"ssemu/internal/resources"
	"ssemu/internal/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CharacterCreateReqHandler struct{}

const CharacterCreateReq network.PacketId = 0x09
const CharacterCreateAck network.PacketId = 0x0A

func (h *CharacterCreateReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	slot := p.ReadU8()
	charMask := domain.CharacterMask(p.ReadU32())

	span.SetAttributes(
		attribute.Int("slot", int(slot)),
		attribute.Int("charMask", int(charMask)),
	)

	if !isCharMaskValid(charMask) {
		return utils.NewValidationError("invalid charmask")
	}

	ack := network.NewPacket(CharacterCreateAck)
	defer ack.Free()

	characters, err := dbGetPlayerCharacters(ctx, s.UserId())
	if err != nil {
		return err
	}

	if len(characters) >= 3 {
		return utils.NewValidationError("cannot create more than 3 characters")
	}

	for _, char := range characters {
		if slot == char.Slot {
			return utils.NewValidationError("create character slot conflicted")
		}
	}

	active := len(characters) == 0

	if err := dbCreateCharacter(ctx, s.UserId(), slot, charMask, active); err != nil {
		return err
	}

	ack.WriteU8(slot)
	ack.WriteU32(uint32(charMask))
	ack.WriteU8(1) // skill count
	ack.WriteU8(3) // weapon count
	s.Send(&ack)

	span.SetAttributes(attribute.Bool("created", true))

	s.Logger().Debug("user created a character")

	return nil
}

func isCharMaskValid(mask domain.CharacterMask) bool {
	var hair, face, shirt, pants byte
	for _, item := range resources.GetDefaultItems(mask.Gender()) {
		if item.Value == "Hair" {
			hair++
		} else if item.Value == "Face" {
			face++
		} else if item.Value == "Coat" {
			shirt++
		} else if item.Value == "Pants" {
			pants++
		}
	}
	valid := mask.Hair() <= hair &&
		mask.Face() <= face &&
		mask.Shirt() <= shirt &&
		mask.Pants() <= pants
	return valid
}

func dbCreateCharacter(ctx context.Context, uid domain.UserId, slot byte, mask domain.CharacterMask, active bool) error {
	ctx, span := database.GetTracer().Start(ctx, "dbCreateCharacter")
	defer span.End()
	const cmd string = `
		INSERT INTO player_characters (
			user_id, slot, mask, weapon1, weapon2, weapon3, skill, 
			hair, face, shirt, pants, shoes, gloves, accessory, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := database.GetConn().ExecContext(ctx, cmd,
		uid, slot, mask, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, active,
	); err != nil {
		return err
	}
	return nil
}

func (h *CharacterCreateReqHandler) AllowAnonymous() bool { return false }

func (h *CharacterCreateReqHandler) GetHandlerName() string { return "CharacterCreateReqHandler" }
