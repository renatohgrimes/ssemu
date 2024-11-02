package resources

import "encoding/xml"

type x4ItemsRoot struct {
	XMLName    xml.Name         `xml:"s4_items"`
	Categories []x4ItemCategory `xml:"category"`
}

type x4ItemCategory struct {
	Id            int                 `xml:"id,attr"`
	SubCategories []x4ItemSubCategory `xml:"sub_category"`
}

type x4ItemSubCategory struct {
	Id    int      `xml:"id,attr"`
	Items []x4Item `xml:"item"`
}

type x4Item struct {
	Number   int             `xml:"number,attr"`
	Products []x4ItemProduct `xml:"product"`
	Base     x4ItemBase      `xml:"base"`
	Costume  x4ItemCostume   `xml:"costume"`
}

type x4ItemProduct struct {
	Id                     int    `xml:"id,attr"`
	TermContract           string `xml:"term_contract,attr"`
	PenPrice               int    `xml:"gm_price,attr"`
	CashPrice              int    `xml:"cash_price,attr"`
	DurabilityInitialValue int    `xml:"durability_inital_value,attr"`
	RemainingSecond        int    `xml:"remaining_second,attr"`
	RefundEnable           bool   `xml:"refund_enable,attr"`
	Development            bool   `xml:"dev,attr"`
}

type x4ItemBase struct {
	RequireLevel    int    `xml:"require_level,attr"`
	ExpBoostPercent int    `xml:"exp_boost_percent,attr"`
	RequireLicense  string `xml:"require_license,attr"`
	EffectGroup     int    `xml:"effect_group,attr"`
}

type x4ItemCostume struct {
	WearingPossibleSex string `xml:"wearing_possible_sex,attr"`
}
