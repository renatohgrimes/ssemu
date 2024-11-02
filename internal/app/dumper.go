package app

import (
	"context"
	"ssemu/internal/network"
)

// development use only
type DumperHandler struct{}

func (h *DumperHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	p.DumpPacket("DumperHandler")
	return nil
}

func (h *DumperHandler) AllowAnonymous() bool { return false }

func (h *DumperHandler) GetHandlerName() string { return "DumperHandler" }
