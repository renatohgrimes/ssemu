package handlers

import (
	"context"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/network"
)

type RoomLoadedReqHandler struct {
}

const RoomLoadedReq network.PacketId = 0x81

func (h *RoomLoadedReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	channel, err := channels.GetUserChannel(s.UserId())
	if err != nil {
		return err
	}
	if err := channel.HandlePlayerRoomLoaded(ctx, s); err != nil {
		return err
	}
	return nil
}

func (h *RoomLoadedReqHandler) AllowAnonymous() bool { return false }

func (h *RoomLoadedReqHandler) GetHandlerName() string { return "RoomLoadedReqHandler" }
