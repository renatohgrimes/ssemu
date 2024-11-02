package domain

type CharacterMask uint32

func (c CharacterMask) Gender() CharacterGender {
	return CharacterGender(c >> 0 & 1)
}

func (c CharacterMask) Hair() byte {
	return byte(c >> 1 & 3)
}

func (c CharacterMask) Face() byte {
	return byte(c >> 7 & 1)
}

func (c CharacterMask) Shirt() byte {
	return byte(c >> 13 & 3)
}

func (c CharacterMask) Pants() byte {
	return byte(c >> 23 & 3)
}
