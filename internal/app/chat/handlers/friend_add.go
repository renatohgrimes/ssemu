package handlers

import (
	"context"
	"ssemu/internal/network"
)

type FriendAddReqHandler struct{}

const FriendAddReq network.PacketId = 0x1A
const FriendAddAck network.PacketId = 0x1B

func (h *FriendAddReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	ack := network.NewPacket(BlockAddAck)
	defer ack.Free()
	ack.WriteU64(0)
	ack.WriteU8(2)
	s.Send(&ack)
	return nil
}

func (h *FriendAddReqHandler) AllowAnonymous() bool { return false }

func (h *FriendAddReqHandler) GetHandlerName() string { return "FriendAddReqHandler" }
