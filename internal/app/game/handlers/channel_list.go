package handlers

import (
	"context"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/network"
)

type ChannelListReqHandler struct{}

const ChannelListReq network.PacketId = 0x2F
const ChannelListAck network.PacketId = 0x2A
const RoomListAck network.PacketId = 0x30

type ChannelListType byte

const (
	Room1   ChannelListType = 3
	Room2   ChannelListType = 4
	Channel ChannelListType = 5
)

func (h *ChannelListReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	listType := ChannelListType(p.ReadU8())
	if listType == Channel {
		ack := network.NewPacket(ChannelListAck)
		defer ack.Free()
		channels := channels.List()
		ack.WriteU16(uint16(len(channels)))
		for _, channel := range channels {
			ack.WriteU16(uint16(channel.GetId()))
			ack.WriteU16(uint16(channel.GetSessionCount()))
		}
		s.Send(&ack)
	} else if listType == Room1 || listType == Room2 {
		channel, err := channels.GetUserChannel(s.UserId())
		if err != nil {
			return err
		}
		ack := network.NewPacket(RoomListAck)
		defer ack.Free()
		ack.WriteU16(uint16(channel.GetRoomCount()))
		for room := range channel.EnumerateRooms() {
			ack.WriteU32(room.GetId())
			ack.WriteU8(room.GetConnectingCount())
			ack.WriteU8(byte(room.GetPlayerCount()))
			ack.WriteU8(byte(room.GetState()))
			ack.WriteU8(100) // ping
			ack.WriteU32(uint32(room.GetSettings()))
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
		}
		s.Send(&ack)
	}
	return nil
}

func (h *ChannelListReqHandler) AllowAnonymous() bool { return false }

func (h *ChannelListReqHandler) GetHandlerName() string { return "ChannelListReqHandler" }
