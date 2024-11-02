package handlers

import (
	"context"
	"errors"
	"fmt"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"ssemu/internal/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UseItemReqHandler struct{}

const UseItemReq network.PacketId = 0x15
const UseItemAck network.PacketId = 0x16

func (h *UseItemReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	cmd := p.ReadU8()
	charSlot := p.ReadU8()
	weaponSlot := p.ReadU8()
	itemId := domain.PlayerItemId(p.ReadU64())

	playerInventory := domain.Mocks.Inventory

	activeCharacter, err := dbGetActiveCharacter(ctx, s.UserId())
	if err != nil {
		return err
	}

	if charSlot != activeCharacter.Slot {
		return utils.NewValidationError("user trying to equip item to a invalid character")
	}

	item, err := inventoryGetItem(itemId, playerInventory)
	if err != nil {
		return err
	}

	tmp := itemId
	unequipping := cmd != 1
	if unequipping {
		tmp = 0
	}

	if unequipping && cmd == 1 {
		return utils.NewValidationError("unexpected command")
	}

	span.SetAttributes(
		attribute.Bool("unequipping", unequipping),
		attribute.Int64("itemId", int64(itemId)),
	)

	if item.Category == 1 { // clothes
		if item.SubCategory == 0 { // hair
			activeCharacter.Hair = tmp
		} else if item.SubCategory == 1 { // face
			activeCharacter.Face = tmp
		} else if item.SubCategory == 2 { // shirt
			activeCharacter.Shirt = tmp
		} else if item.SubCategory == 3 { // pants
			activeCharacter.Pants = tmp
		} else if item.SubCategory == 4 { // gloves
			activeCharacter.Gloves = tmp
		} else if item.SubCategory == 5 { // shoes
			activeCharacter.Shoes = tmp
		} else if item.SubCategory == 6 { // accessory
			activeCharacter.Accessory = tmp
		} else {
			return utils.NewValidationError("invalid clothes subcategory")
		}
	} else if item.Category == 2 { // weapons
		if !unequipping && (activeCharacter.Weapon1 == tmp || activeCharacter.Weapon2 == tmp || activeCharacter.Weapon3 == tmp) {
			return utils.NewValidationError("cannot overlap weapons")
		}
		if weaponSlot == 0 {
			activeCharacter.Weapon1 = tmp
		} else if weaponSlot == 1 {
			activeCharacter.Weapon2 = tmp
		} else if weaponSlot == 2 {
			activeCharacter.Weapon3 = tmp
		} else {
			return utils.NewValidationError("invalid weapon slot")
		}
	} else if item.Category == 3 { // skills
		activeCharacter.Skill = tmp
	} else {
		return utils.NewValidationError("invalid category")
	}

	if err := dbUpdateCharacter(ctx, s.UserId(), activeCharacter); err != nil {
		return err
	}

	ack := network.NewPacket(UseItemAck)
	defer ack.Free()
	ack.WriteU8(cmd)
	ack.WriteU8(charSlot)
	ack.WriteU8(weaponSlot)
	ack.WriteU64(uint64(itemId))

	s.Send(&ack)

	return nil
}

func inventoryGetItem(itemId domain.PlayerItemId, inventory []domain.PlayerItem) (domain.PlayerItem, error) {
	for _, item := range inventory {
		if itemId == item.Id {
			return item, nil
		}
	}
	return domain.PlayerItem{}, fmt.Errorf("item id %d not found in inventory", itemId)
}

func dbGetActiveCharacter(ctx context.Context, userId domain.UserId) (domain.PlayerCharacter, error) {
	ctx, span := database.GetTracer().Start(ctx, "dbGetActiveCharacter")
	defer span.End()
	_, activeSlot, err := dbGetCharacterSlotData(ctx, userId)
	if err != nil {
		return domain.PlayerCharacter{}, err
	}
	characters, err := dbGetPlayerCharacters(ctx, userId)
	if err != nil {
		return domain.PlayerCharacter{}, err
	}
	for _, character := range characters {
		if character.Slot == activeSlot {
			return character, nil
		}
	}
	return domain.PlayerCharacter{}, errors.New("active character not found")
}

func dbUpdateCharacter(ctx context.Context, userId domain.UserId, chara domain.PlayerCharacter) error {
	ctx, span := database.GetTracer().Start(ctx, "dbUpdateCharacter")
	defer span.End()
	const query string = `
		UPDATE player_characters SET
			mask = ?, weapon1 = ?, weapon2 = ?, weapon3 = ?, skill = ?, 
			hair = ?, face = ?, shirt = ?, pants = ?, shoes = ?, gloves = ?, accessory = ?
		WHERE
			user_id = ? AND slot = ?`
	if _, err := database.GetConn().ExecContext(ctx, query,
		chara.Mask, chara.Weapon1, chara.Weapon2, chara.Weapon3, chara.Skill,
		chara.Hair, chara.Face, chara.Shirt, chara.Pants, chara.Shoes, chara.Gloves, chara.Accessory,
		userId, chara.Slot,
	); err != nil {
		return err
	}
	return nil
}

func (h *UseItemReqHandler) AllowAnonymous() bool { return false }

func (h *UseItemReqHandler) GetHandlerName() string { return "UseItemReqHandler" }
