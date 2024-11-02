package channels

import (
	"context"
	"errors"
	"fmt"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"ssemu/internal/resources"
	"sync"
	"time"
)

var ErrChannelCapacityExceeded error = errors.New("channel capacity exceeded")

type Channel struct {
	id       uint32
	name     string
	sessions map[domain.UserId]network.Session
	rooms    map[uint32]*room
	mutex    *sync.Mutex
}

func NewChannel(id uint32, name string) Channel {
	return Channel{
		id:       id,
		name:     name,
		mutex:    &sync.Mutex{},
		sessions: make(map[domain.UserId]network.Session),
		rooms:    make(map[uint32]*room),
	}
}

func (c *Channel) Join(session network.Session) error {
	if _, exists := c.sessions[session.UserId()]; exists {
		return nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.GetSessionCount() >= c.GetCapacity() {
		return ErrChannelCapacityExceeded
	}
	c.sessions[session.UserId()] = session
	session.Logger().Debug("added to server channels")
	return nil
}

func (c *Channel) Leave(userId domain.UserId) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if session, exists := c.sessions[userId]; exists {
		if room, err := c.getUserRoom(userId); err == nil {
			room.removePlayer(userId, None)
			if len(room.players) == 0 {
				c.removeRoom(room.id)
			}
		}
		delete(c.sessions, userId)
		session.Logger().Debug("removed from server channels")
	}
}

func (c *Channel) CreateRoom(name string, password uint32, settings RoomSettings, timeLimit time.Duration, scoreLimit byte, noIntrusion bool) *room {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var counter uint32 = 1
	for ; ; counter++ {
		if _, exists := c.rooms[counter]; exists {
			continue
		}
		break
	}
	room := NewRoom(counter, name, password, settings, timeLimit, scoreLimit, noIntrusion)
	c.rooms[counter] = &room
	return &room
}

func (c *Channel) JoinRoom(ctx context.Context, session network.Session, roomId uint32, password uint32) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	room, exists := c.rooms[roomId]
	if !exists {
		return fmt.Errorf("room %d not found", roomId)
	}
	if err := room.join(ctx, session, password); err != nil {
		return err
	}
	return nil
}

func (c *Channel) LeaveRoom(userId domain.UserId) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if room, err := c.getUserRoom(userId); err == nil {
		room.removePlayer(userId, None)
		if len(room.players) == 0 {
			c.removeRoom(room.id)
		}
	}
}

func (c *Channel) HandlePlayerRoomLoaded(ctx context.Context, session network.Session) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	room, err := c.getUserRoom(session.UserId())
	if err != nil {
		return err
	}
	err = room.handlePlayerRoomLoaded(ctx, session)
	return err
}

func (c *Channel) removeRoom(roomId uint32) {
	if room, exists := c.rooms[roomId]; exists {
		for _, player := range room.players {
			room.removePlayer(player.session.UserId(), None)
		}
		delete(c.rooms, roomId)
	}
}

func (c *Channel) getUserRoom(userId domain.UserId) (*room, error) {
	for _, room := range c.rooms {
		if _, exists := room.players[userId]; exists {
			return room, nil
		}
	}
	return nil, fmt.Errorf("user %d room not found", userId)
}

func (c *Channel) EnumerateSessions() func(func(network.Session) bool) {
	return func(yield func(network.Session) bool) {
		for _, session := range c.sessions {
			yield(session)
		}
	}
}

func (c *Channel) EnumerateRooms() func(func(*room) bool) {
	return func(yield func(*room) bool) {
		for _, room := range c.rooms {
			yield(room)
		}
	}
}

func (c *Channel) GetId() uint32 { return c.id }

func (c *Channel) GetName() string { return c.name }

func (c *Channel) GetSessionCount() int { return len(c.sessions) }

func (c *Channel) GetCapacity() int { return resources.GetChannelCapacity() }

func (c *Channel) GetRoomCount() int { return len(c.rooms) }
