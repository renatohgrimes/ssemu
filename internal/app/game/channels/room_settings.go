package channels

type RoomSettings uint32

func (r RoomSettings) PlayerCount() byte {
	cmd := byte(r >> 16)
	if cmd == 8 {
		return 12
	} else if cmd == 7 {
		return 10
	} else if cmd == 6 {
		return 8
	} else if cmd == 5 {
		return 6
	} else if cmd == 3 {
		return 4
	}
	return 12
}

func (r RoomSettings) GameMode() GameMode { return GameMode(byte(r) >> 4) }

func (r RoomSettings) MapId() byte { return byte(r >> 8) }

func (r RoomSettings) SpectatorCount() byte { return byte(r >> 24) }

func (r RoomSettings) HasPassword() byte { return byte(r >> 1 & 1) }
