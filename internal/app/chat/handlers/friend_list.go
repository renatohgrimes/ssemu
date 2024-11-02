package handlers

import (
	"context"
	"ssemu/internal/network"
)

type FriendListReqHandler struct{}

const FriendListReq network.PacketId = 0x24
const FriendListAck network.PacketId = 0x25

func (h *FriendListReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *FriendListReqHandler) AllowAnonymous() bool { return false }

func (h *FriendListReqHandler) GetHandlerName() string { return "FriendListReqHandler" }
