package chat

import (
	"errors"
	"ssemu/internal/app/chat/handlers"
	"ssemu/internal/network"
)

func Setup(chatServer network.Server, gameServer network.Server) error {
	return errors.Join(
		chatServer.RegisterHandler(handlers.LoginReq, &handlers.LoginReqHandler{
			GameServer: gameServer,
			ChatServer: chatServer,
		}),
		chatServer.RegisterHandler(handlers.KeepAliveReq, &handlers.KeepAliveReqHandler{}),
		chatServer.RegisterHandler(handlers.MessageReq, &handlers.MessageReqHandler{}),
		chatServer.RegisterHandler(handlers.WhisperReq, &handlers.WhisperReqHandler{}),
		chatServer.RegisterHandler(handlers.BlockAddReq, &handlers.BlockAddReqHandler{}),
		chatServer.RegisterHandler(handlers.BlockListReq, &handlers.BlockListReqHandler{}),
		chatServer.RegisterHandler(handlers.ChannelEnterReq, &handlers.ChannelEnterReqHandler{}),
		chatServer.RegisterHandler(handlers.ChannelLeaveReq, &handlers.ChannelLeaveReqHandler{}),
		chatServer.RegisterHandler(handlers.ChannelListReq, &handlers.ChannelListReqHandler{}),
		chatServer.RegisterHandler(handlers.GetChatSettingsReq, &handlers.GetChatSettingsReqHandler{}),
		chatServer.RegisterHandler(handlers.SetChatSettingsReq, &handlers.SetChatSettingsReqHandler{}),
		chatServer.RegisterHandler(handlers.SetChatStateReq, &handlers.SetChatStateReqHandler{}),
		chatServer.RegisterHandler(handlers.FriendAddReq, &handlers.FriendAddReqHandler{}),
		chatServer.RegisterHandler(handlers.FriendListReq, &handlers.FriendListReqHandler{}),
		chatServer.RegisterHandler(handlers.CombiListReq, &handlers.CombiListReqHandler{}),
	)
}
