package handlers

import (
	"context"
	"ssemu/internal/network"
)

type BlockAddReqHandler struct{}

const BlockAddReq network.PacketId = 0x53
const BlockAddAck network.PacketId = 0x54

func (h *BlockAddReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	ack := network.NewPacket(BlockAddAck)
	defer ack.Free()
	ack.WriteU8(2)
	s.Send(&ack)
	return nil
}

func (h *BlockAddReqHandler) AllowAnonymous() bool { return false }

func (h *BlockAddReqHandler) GetHandlerName() string { return "BlockAddReqHandler" }
