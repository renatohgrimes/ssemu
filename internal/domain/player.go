package domain

import (
	"time"
)

type Player struct {
	UserId         UserId
	Nickname       Nickname
	CreatedUtc     time.Time
	TutorialStatus PlayerTutorialStatus
}

func NewPlayer(userId UserId, nickname Nickname) Player {
	return Player{
		UserId:         userId,
		Nickname:       nickname,
		CreatedUtc:     time.Now().UTC(),
		TutorialStatus: TutorialPending,
	}
}

type PlayerCharacter struct {
	Slot      byte
	Mask      CharacterMask
	Weapon1   PlayerItemId
	Weapon2   PlayerItemId
	Weapon3   PlayerItemId
	Skill     PlayerItemId
	Hair      PlayerItemId
	Face      PlayerItemId
	Shirt     PlayerItemId
	Pants     PlayerItemId
	Shoes     PlayerItemId
	Gloves    PlayerItemId
	Accessory PlayerItemId
}

type PlayerItemId uint64

type PlayerItem struct {
	Id           PlayerItemId
	Category     byte
	SubCategory  byte
	Number       uint16
	Product      byte
	EffectGroup  uint32
	SellPrice    uint32
	PurchaseTime int64
	ExpireTime   int64
	Energy       int32
	TimeLeft     int32
}

type PlayerTutorialStatus uint32

const (
	TutorialPending PlayerTutorialStatus = 0
	TutorialDone    PlayerTutorialStatus = 3
)
