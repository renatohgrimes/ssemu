package game

import (
	"errors"
	"ssemu/internal/app/game/channels"
	"ssemu/internal/app/game/handlers"
	"ssemu/internal/network"
)

func Setup(gameServer network.Server, authServer network.Server) error {
	channels.Load()
	return errors.Join(
		gameServer.RegisterHandler(handlers.LoginReq, &handlers.LoginReqHandler{
			AuthServer: authServer,
			GameServer: gameServer,
		}),
		gameServer.RegisterHandler(handlers.KeepAliveReq, &handlers.KeepAliveReqHandler{}),
		gameServer.RegisterHandler(handlers.TimeSyncReq, &handlers.TimeSyncReqHandler{}),
		gameServer.RegisterHandler(handlers.GetNicknameAvailabilityReq, &handlers.GetNicknameAvailabilityReqHandler{}),
		gameServer.RegisterHandler(handlers.AdminShowWindowReq, &handlers.AdminShowWindowReqHandler{}),
		gameServer.RegisterHandler(handlers.LogoutReq, &handlers.LogoutReqHandler{}),
		gameServer.RegisterHandler(handlers.CreatePlayerReq, &handlers.CreatePlayerReqHandler{}),
		gameServer.RegisterHandler(handlers.NatInfoReq, &handlers.NatInfoReqHandler{}),
		gameServer.RegisterHandler(handlers.CharacterCreateReq, &handlers.CharacterCreateReqHandler{}),
		gameServer.RegisterHandler(handlers.CharacterSelectReq, &handlers.CharacterSelectReqHandler{}),
		gameServer.RegisterHandler(handlers.ChannelListReq, &handlers.ChannelListReqHandler{}),
		gameServer.RegisterHandler(handlers.NewEventReq, &handlers.NewEventReqHandler{}),
		gameServer.RegisterHandler(handlers.TutorialDoneReq, &handlers.TutorialDoneReqHandler{}),
		gameServer.RegisterHandler(handlers.ChannelEnterReq, &handlers.ChannelEnterReqHandler{}),
		gameServer.RegisterHandler(handlers.RefreshEquipmentReq, &handlers.RefreshEquipmentReqHandler{}),
		gameServer.RegisterHandler(handlers.RefreshInventoryReq, &handlers.RefreshInventoryReqHandler{}),
		gameServer.RegisterHandler(handlers.ChannelLeaveReq, &handlers.ChannelLeaveReqHandler{}),
		gameServer.RegisterHandler(handlers.AdminCommandReq, &handlers.AdminCommandReqHandler{}),
		gameServer.RegisterHandler(handlers.RefundItemReq, &handlers.RefundItemReqHandler{}),
		gameServer.RegisterHandler(handlers.UseItemReq, &handlers.UseItemReqHandler{}),
		gameServer.RegisterHandler(handlers.RoomQuickJoinReq, &handlers.RoomQuickJoinReqHandler{}),
		gameServer.RegisterHandler(handlers.ItemBuyReq, &handlers.ItemBuyReqHandler{}),
		gameServer.RegisterHandler(handlers.RepairItemReq, &handlers.RepairItemReqHandler{}),
		gameServer.RegisterHandler(handlers.RepairEquippedItemReq, &handlers.RepairEquippedItemReqHandler{}),
		gameServer.RegisterHandler(handlers.CharacterDeleteReq, &handlers.CharacterDeleteReqHandler{}),
		gameServer.RegisterHandler(handlers.RoomCreateReq, &handlers.RoomCreateReqHandler{}),
		gameServer.RegisterHandler(handlers.TunnelJoinReq, &handlers.TunnelJoinReqHandler{}),
		gameServer.RegisterHandler(handlers.RoomLoadedReq, &handlers.RoomLoadedReqHandler{}),
		gameServer.RegisterHandler(handlers.RoomJoinReq, &handlers.RoomJoinReqHandler{}),
	)
}
