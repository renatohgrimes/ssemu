package domain

type CharacterGender byte

const (
	Male   CharacterGender = 0
	Female CharacterGender = 1
)

type LicenseId byte

const (
	None            LicenseId = 0
	PlasmaSword     LicenseId = 1
	CounterSword    LicenseId = 2
	StormBat        LicenseId = 26
	SubmachineGun   LicenseId = 3
	Revolver        LicenseId = 4
	SemiRifle       LicenseId = 25
	HeavymachineGun LicenseId = 5
	RailGun         LicenseId = 6
	Cannonade       LicenseId = 7
	Sentrygun       LicenseId = 8
	MineGun         LicenseId = 10
	MindShock       LicenseId = 12
	Anchoring       LicenseId = 13
	Flying          LicenseId = 14
	Invisible       LicenseId = 15
	Shield          LicenseId = 17
	Block           LicenseId = 18
	Bind            LicenseId = 19
)
