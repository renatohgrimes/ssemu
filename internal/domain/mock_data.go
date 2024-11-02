package domain

type MockData struct {
	Inventory []PlayerItem
	Pen       uint32
	Cash      uint32
	Level     uint8
	Exp       uint32
}

var Mocks = MockData{
	Pen:   9999999,
	Cash:  9999999,
	Level: 100,
	Exp:   63703100,
	Inventory: []PlayerItem{
		mockPlayerItem(1, 3, 10, 1),
		mockPlayerItem(5, 2, 0, 1),
	},
}

func mockPlayerItem(id uint64, category byte, subcategory byte, number uint16) PlayerItem {
	const maxEnergyEquip = 2400
	const maxEnergyClothes = 900
	var energy int32 = maxEnergyEquip
	if category == 1 {
		energy = maxEnergyClothes
	}
	return PlayerItem{
		Id:           PlayerItemId(id),
		Category:     category,
		SubCategory:  subcategory,
		Number:       number,
		Product:      1,
		ExpireTime:   -1,
		Energy:       energy,
		TimeLeft:     -1,
		EffectGroup:  0,
		SellPrice:    0,
		PurchaseTime: 1500000000,
	}
}
