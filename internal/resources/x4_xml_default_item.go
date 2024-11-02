package resources

import "encoding/xml"

type x4DefaultItemRoot struct {
	XMLName xml.Name            `xml:"s4_default_item"`
	Male    x4DefaultItemMale   `xml:"male"`
	Female  x4DefaultItemFemale `xml:"female"`
}

type x4DefaultItemMale struct {
	XMLName xml.Name        `xml:"male"`
	Items   []x4DefaultItem `xml:"item"`
}

type x4DefaultItemFemale struct {
	XMLName xml.Name        `xml:"female"`
	Items   []x4DefaultItem `xml:"item"`
}

type x4DefaultItem struct {
	XMLName     xml.Name `xml:"item"`
	Category    int      `xml:"category,attr"`
	SubCategory int      `xml:"sub_category,attr"`
	Number      int      `xml:"number,attr"`
	Value       string   `xml:",chardata"`
}
