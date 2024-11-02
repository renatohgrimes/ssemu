package relay

import (
	"errors"
	"ssemu/internal/app/relay/handlers"
	"ssemu/internal/network"
)

func Setup(relayServer network.Server, authServer network.Server, gameServer network.Server) error {
	return errors.Join(
		relayServer.RegisterHandler(handlers.LoginReq, &handlers.LoginReqHandler{
			RelayServer: relayServer,
			AuthServer:  authServer,
			GameServer:  gameServer,
		}),
		relayServer.RegisterHandler(handlers.KeepAliveReq, &handlers.KeepAliveReqHandler{}),
		relayServer.RegisterHandler(handlers.DetourReq, &handlers.DetourReqHandler{
			RelayServer: relayServer,
		}),
		relayServer.RegisterHandler(handlers.TunnelJoinReq, &handlers.TunnelJoinReqHandler{}),
		relayServer.RegisterHandler(handlers.TunnelLeaveReq, &handlers.TunnelLeaveReqHandler{}),
		relayServer.RegisterHandler(handlers.TunnelUseReq, &handlers.TunnelUseReqHandler{}),
	)
}
