package handlers

import (
	"context"
	"ssemu/internal/app"
	"ssemu/internal/network"
)

type ServerListReqHandler struct {
	GameServer network.Server
}

const ServerListReq network.PacketId = 0x06
const ServerListAck network.PacketId = 0x07

func (h *ServerListReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	ack := network.NewPacket(ServerListAck)
	defer ack.Free()
	ack.WriteU8(4) // server count
	h.ackWriteServerData(&ack, 1, app.GameServerSettings)
	h.ackWriteServerData(&ack, 1, app.ChatServerSettings)
	h.ackWriteServerData(&ack, 1, app.Nat1ServerSettings)
	h.ackWriteServerData(&ack, 1, app.RelayServerSettings)
	s.Send(&ack)
	return nil
}

//nolint:unparam
func (h *ServerListReqHandler) ackWriteServerData(ack *network.Packet, serverGroup uint16, s network.ServerSettings) {
	ack.WriteU16(serverGroup)
	ack.WriteU8(byte(s.Type))
	ack.WriteStringSlice("ssemu", 40)
	ack.WriteU16(uint16(h.GameServer.GetSessionCount()))
	ack.WriteU16(uint16(app.MaxCapacity))
	ack.WriteU32(h.GameServer.GetPublicEndpoint().GetAddress())
	ack.WriteU16(s.Port)
}

func (h *ServerListReqHandler) AllowAnonymous() bool { return false }

func (h *ServerListReqHandler) GetHandlerName() string { return "ServerListReqHandler" }
