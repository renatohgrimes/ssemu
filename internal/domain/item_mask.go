package domain

type ItemMask uint64

func NewItemMask(category int, subcategory int, number int) ItemMask {
	return ItemMask((category * 1000000) + (subcategory * 10000) + number)
}

func (im ItemMask) Category() int {
	return int(im / 1000000)
}

func (im ItemMask) SubCategory() int {
	remainder := im % 1000000
	subcategory := remainder / 10000
	return int(subcategory)
}

func (im ItemMask) Number() int {
	return int(im % 10000)
}
