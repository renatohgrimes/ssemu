package handlers

import (
	"context"
	"ssemu/internal/network"
)

type BlockListReqHandler struct{}

const BlockListReq network.PacketId = 0x57
const BlockListAck network.PacketId = 0x58

func (h *BlockListReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return nil
}

func (h *BlockListReqHandler) AllowAnonymous() bool { return false }

func (h *BlockListReqHandler) GetHandlerName() string { return "BlockListReqHandler" }
