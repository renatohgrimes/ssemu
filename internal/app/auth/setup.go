package auth

import (
	"errors"
	"ssemu/internal/app/auth/handlers"
	"ssemu/internal/network"
)

func Setup(authServer network.Server, gameServer network.Server) error {
	return errors.Join(
		authServer.RegisterHandler(handlers.LoginReq, &handlers.LoginReqHandler{
			AuthServer: authServer,
		}),
		authServer.RegisterHandler(handlers.ServerListReq, &handlers.ServerListReqHandler{
			GameServer: gameServer,
		}),
	)
}
