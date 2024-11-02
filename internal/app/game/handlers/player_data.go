package handlers

import (
	"context"
	"database/sql"
	"errors"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/trace"
)

const LicenseDataAck network.PacketId = 0x57
const CharacterSlotDataAck network.PacketId = 0x49
const CharacterDataAck network.PacketId = 0x06
const CharacterEquipDataAck network.PacketId = 0x07
const InventoryDataAck network.PacketId = 0x08
const StatsDataAck network.PacketId = 0x05

func sendPlayerData(ctx context.Context, s network.Session) error {
	span := trace.SpanFromContext(ctx)

	// Licenses

	licenses, err := dbGetPlayerLicenses(ctx, s.UserId())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	licenseAck := network.NewPacket(LicenseDataAck)
	defer licenseAck.Free()

	licenseCount := len(licenses)
	licenseAck.WriteU8(byte(licenseCount))

	for _, license := range licenses {
		licenseAck.WriteU8(byte(license))
	}

	s.Send(&licenseAck)

	// Characters

	characterCount, activeSlot, err := dbGetCharacterSlotData(ctx, s.UserId())
	if err != nil {
		return err
	}

	charSlotAck := network.NewPacket(CharacterSlotDataAck)
	defer charSlotAck.Free()

	charSlotAck.WriteU8(characterCount)
	charSlotAck.WriteU8(3) // slot count
	charSlotAck.WriteU8(activeSlot)

	s.Send(&charSlotAck)

	if characterCount > 0 {
		characters, err := dbGetPlayerCharacters(ctx, s.UserId())
		if err != nil {
			return err
		}

		var i byte
		for i = 0; i < characterCount; i++ {
			character := characters[i]

			charDataAck := network.NewPacket(CharacterDataAck)
			defer charDataAck.Free()

			charDataAck.WriteU8(character.Slot)
			charDataAck.WriteU8(1) // skill count
			charDataAck.WriteU8(3) // weapon count
			charDataAck.WriteU32(uint32(character.Mask))
			s.Send(&charDataAck)

			charEquipAck := network.NewPacket(CharacterEquipDataAck)
			defer charEquipAck.Free()

			charEquipAck.WriteU8(character.Slot)
			charEquipAck.WriteU8(1) // skill count
			charEquipAck.WriteU8(3) // weapon count
			charEquipAck.WriteU8(0) // weapon 1
			charEquipAck.WriteU64(uint64(character.Weapon1))
			charEquipAck.WriteU8(1) // weapon 2
			charEquipAck.WriteU64(uint64(character.Weapon2))
			charEquipAck.WriteU8(2) // weapon 3
			charEquipAck.WriteU64(uint64(character.Weapon3))
			charEquipAck.WriteU8(0) // skill
			charEquipAck.WriteU64(uint64(character.Skill))
			charEquipAck.WriteU8(0) // hair
			charEquipAck.WriteU64(uint64(character.Hair))
			charEquipAck.WriteU8(1) // face
			charEquipAck.WriteU64(uint64(character.Face))
			charEquipAck.WriteU8(2) // shirt
			charEquipAck.WriteU64(uint64(character.Shirt))
			charEquipAck.WriteU8(3) // pants
			charEquipAck.WriteU64(uint64(character.Pants))
			charEquipAck.WriteU8(4) // gloves
			charEquipAck.WriteU64(uint64(character.Gloves))
			charEquipAck.WriteU8(5) // shoes
			charEquipAck.WriteU64(uint64(character.Shoes))
			charEquipAck.WriteU8(6) // accessory
			charEquipAck.WriteU64(uint64(character.Accessory))

			s.Send(&charEquipAck)
		}
	}

	// Inventory

	items := domain.Mocks.Inventory

	inventoryAck := network.NewPacket(InventoryDataAck)
	defer inventoryAck.Free()

	inventoryAck.WriteU32(uint32(len(items)))
	for _, item := range items {
		inventoryAck.WriteU64(uint64(item.Id))
		inventoryAck.WriteU8(item.Category)
		inventoryAck.WriteU8(item.SubCategory)
		inventoryAck.WriteU16(item.Number)
		inventoryAck.WriteU8(item.Product)
		inventoryAck.WriteU32(item.EffectGroup)
		inventoryAck.WriteU32(item.SellPrice)
		inventoryAck.WriteI64(item.PurchaseTime)
		inventoryAck.WriteI64(item.ExpireTime)
		inventoryAck.WriteI32(item.Energy)
		inventoryAck.WriteI32(item.TimeLeft)
	}

	s.Send(&inventoryAck)

	inventorySuccessAck := network.NewPacket(ServerResultAck)
	defer inventorySuccessAck.Free()
	ackWriteServerResult(span, &inventorySuccessAck, InventorySuccess)

	s.Send(&inventorySuccessAck)

	// Account

	isAdmin, err := dbIsUserAdmin(ctx, s.UserId())
	if err != nil {
		return err
	}

	player, err := dbGetPlayer(ctx, s.UserId())
	if err != nil {
		return err
	}

	statsAck := network.NewPacket(StatsDataAck)
	defer statsAck.Free()
	if isAdmin {
		statsAck.WriteU8(1)
	} else {
		statsAck.WriteU8(0)
	}
	statsAck.WriteU8(domain.Mocks.Level)
	statsAck.WriteU32(domain.Mocks.Exp)
	statsAck.WriteU32(0) // unk
	statsAck.WriteU32(uint32(player.TutorialStatus))
	statsAck.WriteStringSlice(string(player.Nickname), 31)
	statsAck.WriteU32(0) // unk
	statsAck.Skip(108)   // playerstats

	s.Send(&statsAck)

	loginSuccessAck := network.NewPacket(ServerResultAck)
	defer loginSuccessAck.Free()
	ackWriteServerResult(span, &loginSuccessAck, LoginSuccess)

	s.Send(&loginSuccessAck)

	return nil
}

func dbGetPlayerLicenses(ctx context.Context, userId domain.UserId) (licenses []domain.LicenseId, err error) {
	ctx, span := database.GetTracer().Start(ctx, "dbGetPlayerLicenses")
	defer span.End()
	const query string = "SELECT license_id FROM player_licenses WHERE user_id = ?"
	db := database.GetConn()
	r, err := db.QueryContext(ctx, query, userId)
	if err != nil {
		return licenses, err
	}
	defer r.Close()
	for r.Next() {
		var id domain.LicenseId
		if err = r.Scan(&id); err != nil {
			return licenses, err
		}
		licenses = append(licenses, id)
	}
	return licenses, nil
}

func dbGetCharacterSlotData(ctx context.Context, userId domain.UserId) (charCount byte, activeSlot byte, err error) {
	ctx, span := database.GetTracer().Start(ctx, "dbGetCharacterSlotData")
	defer span.End()
	const query string = "SELECT is_active, slot FROM player_characters WHERE user_id = ?"
	db := database.GetConn()
	r, err := db.QueryContext(ctx, query, userId)
	if err != nil {
		return 0, 0, err
	}
	defer r.Close()
	for r.Next() {
		var isActive bool
		var slot byte
		if err := r.Scan(&isActive, &slot); err != nil {
			return 0, 0, err
		}
		charCount++
		if isActive {
			activeSlot = slot
		}
	}
	return charCount, activeSlot, nil
}

func dbGetPlayerCharacters(ctx context.Context, userId domain.UserId) (characters []domain.PlayerCharacter, err error) {
	ctx, span := database.GetTracer().Start(ctx, "dbGetPlayerCharacters")
	defer span.End()
	const query string = `
	SELECT 
		slot, mask, weapon1, weapon2, weapon3, skill, 
		hair, face, shirt, pants, shoes, gloves, accessory 
	FROM player_characters WHERE user_id = ?`
	db := database.GetConn()
	r, err := db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for r.Next() {
		var c domain.PlayerCharacter
		if err := r.Scan(
			&c.Slot, &c.Mask, &c.Weapon1, &c.Weapon2, &c.Weapon3, &c.Skill,
			&c.Hair, &c.Face, &c.Shirt, &c.Pants, &c.Shoes, &c.Gloves, &c.Accessory,
		); err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}
	return characters, nil
}

func dbGetPlayer(ctx context.Context, userId domain.UserId) (player domain.Player, err error) {
	ctx, span := database.GetTracer().Start(ctx, "dbGetPlayer")
	defer span.End()
	db := database.GetConn()
	const query string = `SELECT nickname, tutorial_status FROM players WHERE user_id = ?`
	r := db.QueryRowContext(ctx, query, userId)
	if err := r.Scan(&player.Nickname, &player.TutorialStatus); err != nil {
		return player, err
	}
	player.UserId = userId
	return player, nil
}
