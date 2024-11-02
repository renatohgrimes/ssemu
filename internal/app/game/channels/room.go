package channels

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"ssemu/internal/app/game/nat"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"time"
)

var ErrRoomCapacityExceed = errors.New("room capacity exceeded")
var ErrRoomWrongPassword = errors.New("room wrong password")
var ErrRoomCannotEnter = errors.New("room cannot enter")
var ErrRoomChangingRules = errors.New("room changing rules")

const (
	RoomJoinAck          network.PacketId = 0x1A
	RoomInfoAck          network.PacketId = 0x18
	RoomEnterAck         network.PacketId = 0x19
	PlayerLoadedAck      network.PacketId = 0x81
	RoomRefereeChangeAck network.PacketId = 0x8D
	RoomMasterChangeAck  network.PacketId = 0x8C
	RoomLeaveAck         network.PacketId = 0x82
	RoomPlayerLeaveAck   network.PacketId = 0x1B
)

type GameMode byte

const (
	DeathMatch GameMode = 1
	TouchDown  GameMode = 2
	Practice   GameMode = 4
)

type RoomState byte

const (
	MatchMaking RoomState = 1
	Playing     RoomState = 2
	Result      RoomState = 3
)

type RoomTiming byte

const (
	NotPlaying RoomTiming = 0
	FirstHalf  RoomTiming = 1
	HalfTime   RoomTiming = 2
	SecondHalf RoomTiming = 3
)

type RoomPlayerState byte

const (
	Alive      RoomPlayerState = 0
	Dead       RoomPlayerState = 1
	Waiting    RoomPlayerState = 2
	Spectating RoomPlayerState = 3
	Lobby      RoomPlayerState = 4
)

type RoomTeam byte

const (
	Neutral RoomTeam = 0
	Alpha   RoomTeam = 1
	Beta    RoomTeam = 2
)

type RoomPlayerMode byte

const (
	Player    RoomPlayerMode = 1
	Spectator RoomPlayerMode = 2
)

type RoomLeaveReason byte

const (
	None RoomLeaveReason = 0
	Kick RoomLeaveReason = 1
)

type roomPlayer struct {
	session      network.Session
	state        RoomPlayerState
	isConnecting bool
	slotId       byte
	isMaster     bool
	isReferee    bool
	team         RoomTeam
	mode         RoomPlayerMode
	nickname     domain.Nickname
}

type room struct {
	id            uint32
	name          string
	password      uint32
	timeLimit     time.Duration
	scoreLimit    byte
	noIntrusion   bool
	state         RoomState
	timing        RoomTiming
	settings      RoomSettings
	players       map[domain.UserId]*roomPlayer
	scoreAlpha    byte
	scoreBeta     byte
	gameStartUtc  time.Time
	capacity      int
	kickedPlayers []domain.UserId
}

func NewRoom(id uint32, name string, password uint32, settings RoomSettings, timeLimit time.Duration, scoreLimit byte, noIntrusion bool) room {
	capacity := settings.PlayerCount() + settings.SpectatorCount()
	if settings.GameMode() == Practice {
		capacity = 1
	}
	return room{
		id:            id,
		name:          name,
		password:      password,
		timeLimit:     timeLimit,
		scoreLimit:    scoreLimit,
		noIntrusion:   noIntrusion,
		state:         MatchMaking,
		timing:        NotPlaying,
		players:       make(map[domain.UserId]*roomPlayer),
		scoreAlpha:    0,
		scoreBeta:     0,
		gameStartUtc:  time.Time{},
		settings:      settings,
		capacity:      int(capacity),
		kickedPlayers: make([]domain.UserId, 0, 10),
	}
}

func (r *room) join(ctx context.Context, session network.Session, password uint32) error {
	if _, exists := r.players[session.UserId()]; exists {
		return nil
	}

	if _, err := nat.Get(session.UserId()); err != nil {
		return ErrRoomCannotEnter
	}

	if r.password != password {
		return ErrRoomWrongPassword
	}

	if len(r.players)+1 > r.capacity {
		return ErrRoomCapacityExceed
	}

	if r.state == Result {
		return ErrRoomCannotEnter
	}

	if r.isUserKicked(session.UserId()) {
		return ErrRoomCannotEnter
	}

	nickname, err := dbGetNickname(ctx, session.UserId())
	if err != nil {
		return err
	}

	slotId := r.getSlotId()

	player := &roomPlayer{
		session:      session,
		state:        Lobby,
		isConnecting: true,
		slotId:       slotId,
		isMaster:     false,
		isReferee:    false,
		team:         r.getAvailableTeam(),
		mode:         r.getAvailablePlayerMode(),
		nickname:     nickname,
	}

	if len(r.players) == 0 {
		player.isMaster = true
		player.isReferee = true
	}

	r.players[session.UserId()] = player

	joinAck := network.NewPacket(RoomJoinAck)
	defer joinAck.Free()

	joinAck.WriteU32(r.id)
	joinAck.WriteU32(uint32(r.settings))
	joinAck.WriteU32(uint32(r.state))
	joinAck.WriteU32(uint32(r.timing))
	joinAck.WriteU32(uint32(r.timeLimit.Milliseconds()))
	if r.state == Playing {
		gameElapsed := time.Now().UTC().Sub(r.gameStartUtc).Milliseconds()
		joinAck.WriteU32(uint32(gameElapsed))
	} else {
		joinAck.WriteU32(0)
	}
	joinAck.WriteU32(uint32(r.scoreLimit))
	joinAck.WriteU8(0)   // friendly
	joinAck.WriteU8(0)   // balanced
	joinAck.WriteU8(0)   // min level
	joinAck.WriteU8(100) // max level
	joinAck.WriteU8(0)   // equip limit - unlimited
	if r.noIntrusion {
		joinAck.WriteU8(1)
	} else {
		joinAck.WriteU8(0)
	}

	session.Send(&joinAck)

	infoAck := network.NewPacket(RoomInfoAck)
	defer infoAck.Free()

	infoAck.WriteU8(player.slotId)
	infoAck.WriteU32(r.id)
	infoAck.WriteU32(0) // unk

	session.Send(&infoAck)

	enterAck := network.NewPacket(RoomEnterAck)
	defer enterAck.Free()

	enterAck.WriteU8(0) // unk
	enterAck.WriteU8(byte(len(r.players)))
	for _, player := range r.players {
		nickname, err := dbGetNickname(ctx, player.session.UserId())
		if err != nil {
			return err
		}
		info, err := nat.Get(player.session.UserId())
		if err != nil {
			return err
		}
		enterAck.WriteU32(info.PrivateAddress)
		enterAck.WriteU16(info.PrivatePort)
		enterAck.WriteU32(info.PublicAddress)
		enterAck.WriteU16(info.PublicPort)
		enterAck.WriteU16(info.NatUnk)
		enterAck.WriteU8(info.ConnectionType)
		enterAck.WriteU64(uint64(player.session.UserId()))
		enterAck.WriteU8(player.slotId)
		enterAck.WriteU32(0) // unk
		enterAck.WriteU8(1)  // unk
		enterAck.WriteStringSlice(string(nickname), 31)
	}

	r.broadcast(&enterAck)

	player.session.Logger().LogAttrs(ctx, slog.LevelDebug, "room joined",
		slog.Int("room", int(r.id)),
		slog.Bool("isSpectator", player.mode == Spectator),
		slog.Bool("isMaster", player.isMaster),
		slog.Bool("isReferee", player.isReferee),
	)

	return nil
}

func (r *room) handlePlayerRoomLoaded(ctx context.Context, session network.Session) error {
	player, exists := r.players[session.UserId()]
	if !exists {
		return fmt.Errorf("user %d not found in room %d", session.UserId(), r.id)
	}

	if !player.isConnecting {
		return fmt.Errorf("user %d is not in connecting state", session.UserId())
	}

	master := r.getMaster()
	referee := r.getReferee()

	r.changeReferee(referee)
	r.changeMaster(master)

	ack := network.NewPacket(PlayerLoadedAck)
	defer ack.Free()

	ack.WriteU64(uint64(session.UserId()))
	ack.WriteU8(byte(player.team))
	ack.WriteU8(byte(player.mode))
	ack.WriteU32(domain.Mocks.Exp)
	ack.WriteStringSlice(string(player.nickname), 31)

	r.broadcast(&ack)

	player.isConnecting = false

	player.session.Logger().LogAttrs(ctx, slog.LevelDebug, "room loaded",
		slog.Int("room", int(r.id)),
		slog.Int("master", int(master.session.UserId())),
		slog.Int("referee", int(referee.session.UserId())),
	)

	return nil
}

func (r *room) removePlayer(uid domain.UserId, reason RoomLeaveReason) {
	player, exists := r.players[uid]
	if !exists {
		return
	}

	if player.isConnecting {
		return
	}

	if reason == Kick {
		r.kickedPlayers = append(r.kickedPlayers, player.session.UserId())
	}

	if len(r.players) >= 2 {
		if player.isMaster {
			newMaster := r.findNewMaster(player)
			r.changeMaster(newMaster)
		}

		if player.isReferee {
			newReferee := r.findNewReferee(player)
			r.changeReferee(newReferee)
		}
	}

	roomLeaveAck := network.NewPacket(RoomLeaveAck)
	defer roomLeaveAck.Free()

	roomLeaveAck.WriteU64(uint64(uid))
	roomLeaveAck.WriteStringSlice(string(player.nickname), 31)
	roomLeaveAck.WriteU8(byte(reason))

	r.broadcast(&roomLeaveAck)

	playerLeaveAck := network.NewPacket(RoomPlayerLeaveAck)
	defer playerLeaveAck.Free()

	playerLeaveAck.WriteU64(uint64(uid))
	playerLeaveAck.WriteU8(player.slotId)

	r.broadcast(&playerLeaveAck)

	delete(r.players, uid)
	player.session.Logger().Debug("removed from room")
}

func (r *room) changeMaster(player *roomPlayer) {
	currentMaster := r.getMaster()
	if player.session.UserId() != currentMaster.session.UserId() {
		currentMaster.isMaster = false
		player.isMaster = true
	}
	ack := network.NewPacket(RoomMasterChangeAck)
	defer ack.Free()
	ack.WriteU64(uint64(player.session.UserId()))
	r.broadcast(&ack)
	player.session.Logger().Debug("current room master")
}

func (r *room) changeReferee(player *roomPlayer) {
	currentReferee := r.getReferee()
	if player.session.UserId() != currentReferee.session.UserId() {
		currentReferee.isReferee = false
		player.isReferee = true
	}
	ack := network.NewPacket(RoomRefereeChangeAck)
	defer ack.Free()
	ack.WriteU64(uint64(player.session.UserId()))
	r.broadcast(&ack)
	player.session.Logger().Debug("current room referee")
}

func (r *room) broadcast(pkt *network.Packet) {
	for _, player := range r.players {
		player.session.Send(pkt)
	}
}

func (r *room) getSlotId() byte {
	var slotId byte = 2
	for _, player := range r.players {
		if player.slotId == slotId {
			slotId++
			continue
		}
	}
	return slotId
}

func (r *room) getMaster() *roomPlayer {
	for _, player := range r.players {
		if player.isMaster {
			return player
		}
	}
	panic("server error: a room must have a master")
}

func (r *room) getReferee() *roomPlayer {
	for _, player := range r.players {
		if player.isReferee {
			return player
		}
	}
	panic("server error: a room must have a referee")
}

func (r *room) findNewMaster(old *roomPlayer) *roomPlayer {
	for _, player := range r.players {
		if !player.isMaster && old.session.UserId() != player.session.UserId() {
			return player
		}
	}
	panic("server error: new room master not found")
}

func (r *room) findNewReferee(old *roomPlayer) *roomPlayer {
	for _, player := range r.players {
		if !player.isReferee && old.session.UserId() != player.session.UserId() {
			return player
		}
	}
	panic("server error: new room referee not found")
}

func (r *room) getAvailableTeam() RoomTeam {
	alphaCount := r.getTeamCount(Alpha)
	betaCount := r.getTeamCount(Beta)
	if alphaCount < betaCount {
		return Alpha
	} else if betaCount < alphaCount {
		return Beta
	}
	return Alpha
}

func (r *room) getAvailablePlayerMode() RoomPlayerMode {
	if r.getPlayerModeCount(Player) < r.settings.PlayerCount() {
		return Player
	}
	if r.getPlayerModeCount(Spectator) < r.settings.SpectatorCount() {
		return Spectator
	}
	panic("server error: unavailable room slots")
}

func (r *room) getTeamCount(team RoomTeam) byte {
	var count byte = 0
	for _, player := range r.players {
		if team == player.team {
			count++
		}
	}
	return count
}

func (r *room) getPlayerModeCount(mode RoomPlayerMode) byte {
	var count byte = 0
	for _, player := range r.players {
		if mode == player.mode {
			count++
		}
	}
	return count
}

func (r *room) isUserKicked(userId domain.UserId) bool {
	for _, kicked := range r.kickedPlayers {
		if kicked == userId {
			return true
		}
	}
	return false
}

func (r *room) GetConnectingCount() byte {
	var count byte = 0
	for _, player := range r.players {
		if player.isConnecting {
			count++
		}
	}
	return count
}

func (r *room) GetId() uint32 { return r.id }

func (r *room) GetPlayerCount() int { return len(r.players) }

func (r *room) GetState() RoomState { return r.state }

func (r *room) GetSettings() RoomSettings { return r.settings }

func (r *room) GetName() string { return r.name }

func (r *room) GetTimeLimit() time.Duration { return r.timeLimit }

func (r *room) GetScoreLimit() byte { return r.scoreLimit }

func (r *room) IsNoIntrusion() bool { return r.noIntrusion }

func dbGetNickname(ctx context.Context, userId domain.UserId) (nickname domain.Nickname, err error) {
	ctx, span := database.GetTracer().Start(ctx, "dbGetNickname")
	defer span.End()
	db := database.GetConn()
	const query string = `SELECT nickname FROM players WHERE user_id = ? LIMIT 1`
	r := db.QueryRowContext(ctx, query, userId)
	if err := r.Scan(&nickname); err != nil {
		return "", err
	}
	return nickname, nil
}
