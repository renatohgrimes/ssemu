package handlers

import (
	"context"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/network"
	"ssemu/internal/resources"
	"ssemu/internal/utils"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RoomCreateReqHandler struct{}

const RoomCreateReq network.PacketId = 0x11
const RoomCreateAck network.PacketId = 0x32

func (h *RoomCreateReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	channel, err := channels.GetUserChannel(s.UserId())
	if err != nil {
		return err
	}

	name := p.ReadStringSlice(31)
	settings := channels.RoomSettings(p.ReadU32())
	timeLimitMinutes := p.ReadU8()
	scoreLimit := p.ReadU8()
	p.Skip(4)
	password := p.ReadU32()
	p.Skip(5)
	noIntrusion := p.ReadU8() == 1

	span.SetAttributes(
		attribute.String("name", name),
		attribute.Int("settings.GameMode", int(settings.GameMode())),
		attribute.Int("settings.MapId", int(settings.MapId())),
		attribute.Int("settings.PlayerCount", int(settings.PlayerCount())),
		attribute.Int("settings.SpectatorCount", int(settings.SpectatorCount())),
		attribute.Int("timeLimitMinutes", int(timeLimitMinutes)),
		attribute.Int("scoreLimit", int(scoreLimit)),
		attribute.Int("password", int(password)),
		attribute.Bool("noIntrusion", noIntrusion),
	)

	if settings.PlayerCount()+settings.SpectatorCount() > 12 {
		return utils.NewValidationError("cannot overflow 12 player count in a room")
	}

	if settings.GameMode() != channels.DeathMatch && settings.GameMode() != channels.Practice && settings.GameMode() != channels.TouchDown {
		return utils.NewValidationError("invalid game mode")
	}

	if len(name) == 0 {
		return utils.NewValidationError("room name must have at least one character")
	}

	if len(name) > 28 {
		return utils.NewValidationError("room name max length exceeded")
	}

	if settings.PlayerCount() < 4 {
		return utils.NewValidationError("room must have at least 4 players")
	}

	if settings.GameMode() == channels.Practice && settings.SpectatorCount() > 0 {
		return utils.NewValidationError("practice has invalid game mode settings")
	}

	capacity := settings.PlayerCount() + settings.SpectatorCount()
	if settings.SpectatorCount() != 0 && capacity-settings.PlayerCount() != settings.SpectatorCount() {
		return utils.NewValidationError("invalid player to spectator capacity settings")
	}

	character, err := dbGetActiveCharacter(ctx, s.UserId())
	if err != nil {
		return err
	}
	if character.Skill == 0 || (character.Weapon1 == 0 && character.Weapon2 == 0 && character.Weapon3 == 0) {
		return utils.NewValidationError("character must have a skill and at least one weapon")
	}

	if !resources.IsMapValid(int(settings.MapId())) {
		return utils.NewValidationError("map not found")
	}

	if !isGameLimitValid(settings.GameMode(), scoreLimit, timeLimitMinutes) {
		return utils.NewValidationError("room game limit is invalid")
	}

	timeLimit := time.Duration(timeLimitMinutes) * time.Minute

	room := channel.CreateRoom(name, password, settings, timeLimit, scoreLimit, noIntrusion)

	ack := network.NewPacket(RoomCreateAck)
	defer ack.Free()
	ack.WriteU32(room.GetId())
	ack.WriteU32(uint32(room.GetSettings()))
	ack.WriteU32(uint32(room.GetState()))
	ack.WriteU8(100) // ping
	ack.WriteStringSlice(room.GetName(), 31)
	ack.WriteU8(room.GetSettings().HasPassword())
	ack.WriteU32(uint32(room.GetTimeLimit().Milliseconds()))
	ack.WriteU32(uint32(room.GetScoreLimit()))
	ack.WriteU8(0)   // friendly
	ack.WriteU8(0)   // balanced
	ack.WriteU8(0)   // min level
	ack.WriteU8(100) // max level
	ack.WriteU8(0)   // equip limit - unlimited
	if room.IsNoIntrusion() {
		ack.WriteU8(1)
	} else {
		ack.WriteU8(0)
	}

	s.Send(&ack)

	if err := channel.JoinRoom(ctx, s, room.GetId(), password); err != nil {
		return err
	}

	return nil
}

func isGameLimitValid(gameMode channels.GameMode, score byte, timeMinutes byte) bool {
	if gameMode == channels.DeathMatch {
		return resources.IsGameLimitValid(int(score), int(timeMinutes), resources.GetGameScoreLimit().DeathMatch.Scores, resources.GetGameTimeLimit().DeathMatch.Times)
	} else if gameMode == channels.TouchDown {
		return resources.IsGameLimitValid(int(score), int(timeMinutes), resources.GetGameScoreLimit().TouchDown.Scores, resources.GetGameTimeLimit().TouchDown.Times)
	} else if gameMode == channels.Practice {
		return resources.IsGameLimitValid(int(score), int(timeMinutes), resources.GetGameScoreLimit().Practice.Scores, resources.GetGameTimeLimit().Practice.Times)
	}
	return false
}

func (h *RoomCreateReqHandler) AllowAnonymous() bool { return false }

func (h *RoomCreateReqHandler) GetHandlerName() string { return "RoomCreateReqHandler" }
