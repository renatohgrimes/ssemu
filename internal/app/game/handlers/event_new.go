package handlers

import (
	"context"
	"errors"
	"ssemu/internal/network"
)

type NewEventReqHandler struct{}

const NewEventReq network.PacketId = 0x85

func (h *NewEventReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	return errors.New("not implemented")
}

func (h *NewEventReqHandler) AllowAnonymous() bool { return false }

func (h *NewEventReqHandler) GetHandlerName() string { return "NewEventReqHandler" }
