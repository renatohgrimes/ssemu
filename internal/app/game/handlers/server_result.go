package handlers

import (
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const ServerResultAck network.PacketId = 0x02

type ServerResult uint32

const (
	GameServerError          ServerResult = 0
	QuickJoinFailed          ServerResult = 1
	AlreadyPlaying           ServerResult = 2
	NonExistingChannel       ServerResult = 3
	ChannelCapacityExceeded  ServerResult = 4
	ChannelEnter             ServerResult = 5
	RoomCapacityExceed       ServerResult = 6
	RoomChangingRules        ServerResult = 7
	ChannelLeave             ServerResult = 8
	PlayerCannotBeFound      ServerResult = 9
	FailedCreatePlayer       ServerResult = 10
	FailedDeletePlayer       ServerResult = 11
	FailedSelectPlayer       ServerResult = 12
	PlayerCreateSuccess      ServerResult = 13
	NicknameAlreadyUsed      ServerResult = 14
	NicknameAvailable        ServerResult = 15
	PasswordError            ServerResult = 16
	LoginSuccess             ServerResult = 17
	LogSpecifiedIP           ServerResult = 18
	LogForbidden             ServerResult = 19
	UserAlreadyExist         ServerResult = 20
	DBError                  ServerResult = 21
	FailedCreatePlayer2      ServerResult = 22
	CannotEnterChannel       ServerResult = 23
	RequiredChannelLicense   ServerResult = 24
	WearingUnusableItem      ServerResult = 25
	FailResellItemWearing    ServerResult = 26
	FailEnterRoom            ServerResult = 28
	ImpossibleToEnterRoom    ServerResult = 29
	TaskCompensationError    ServerResult = 31
	FailedToRequestTask      ServerResult = 32
	ItemExchangeFailed       ServerResult = 33
	ItemExchangeFailed2      ServerResult = 34
	SelectGameMode           ServerResult = 35
	DBError2                 ServerResult = 36
	InventoryFilledWithItems ServerResult = 37
	LoginFromElsewhere       ServerResult = 38
	InventorySuccess         ServerResult = 39
)

func ackWriteServerResult(span trace.Span, ack *network.Packet, res ServerResult) {
	span.SetAttributes(attribute.Key("result").Int(int(res)))
	ack.WriteU32(uint32(res))
}
