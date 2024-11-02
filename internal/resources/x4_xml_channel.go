package resources

import "encoding/xml"

type x4ChannelRoot struct {
	XMLName  xml.Name         `xml:"channel"`
	Setting  x4ChannelSetting `xml:"setting"`
	Channels []x4ChannelInfo  `xml:"channel_info"`
}

type x4ChannelSetting struct {
	XMLName         xml.Name `xml:"setting"`
	ChannelCapacity int      `xml:"limit_player,attr"`
}

type x4ChannelInfo struct {
	XMLName  xml.Name      `xml:"channel_info"`
	Id       int           `xml:"id,attr"`
	Type     int           `xml:"type,attr"`
	Language x4ChannelLang `xml:"lang"`
}

type x4ChannelLang struct {
	XMLName xml.Name          `xml:"lang"`
	Nations []x4ChannelNation `xml:"nation"`
}

type x4ChannelNation struct {
	XMLName xml.Name `xml:"nation"`
	Code    string   `xml:"name_code,attr"`
}
