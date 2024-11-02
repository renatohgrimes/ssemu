package resources

import "encoding/xml"

type x4GameInfoRoot struct {
	XMLName xml.Name        `xml:"s4_gameInfo"`
	Maps    x4GameInfoMap   `xml:"map"`
	Scores  x4GameInfoScore `xml:"score"`
	Times   x4GameInfoTime  `xml:"time"`
}

type x4GameInfoMap struct {
	XMLName xml.Name            `xml:"map"`
	Values  []x4GameInfoMapData `xml:"data"`
}

type x4GameInfoMapData struct {
	XMLName xml.Name `xml:"data"`
	Id      int      `xml:"id,attr"`
}

type x4GameInfoScore struct {
	XMLName    xml.Name                  `xml:"score"`
	DeathMatch x4GameInfoScoreDeathMatch `xml:"death_match"`
	TouchDown  x4GameInfoScoreTouchDown  `xml:"touch_down"`
	Practice   x4GameInfoScorePractice   `xml:"mission"`
}

type x4GameInfoScoreDeathMatch struct {
	XMLName xml.Name                  `xml:"death_match"`
	Scores  []x4GameInfoScoreRuleData `xml:"data"`
}

type x4GameInfoScoreTouchDown struct {
	XMLName xml.Name                  `xml:"touch_down"`
	Scores  []x4GameInfoScoreRuleData `xml:"data"`
}

type x4GameInfoScorePractice struct {
	XMLName xml.Name                  `xml:"mission"`
	Scores  []x4GameInfoScoreRuleData `xml:"data"`
}

type x4GameInfoScoreRuleData struct {
	XMLName xml.Name `xml:"data"`
	Score   int      `xml:"score,attr"`
}

type x4GameInfoTime struct {
	XMLName    xml.Name                 `xml:"time"`
	DeathMatch x4GameInfoTimeDeathMatch `xml:"death_match"`
	TouchDown  x4GameInfoTimeTouchDown  `xml:"touch_down"`
	Practice   x4GameInfoTimePractice   `xml:"mission"`
}

type x4GameInfoTimeDeathMatch struct {
	XMLName xml.Name                 `xml:"death_match"`
	Times   []x4GameInfoTimeRuleData `xml:"data"`
}

type x4GameInfoTimeTouchDown struct {
	XMLName xml.Name                 `xml:"touch_down"`
	Times   []x4GameInfoTimeRuleData `xml:"data"`
}

type x4GameInfoTimePractice struct {
	XMLName xml.Name                 `xml:"mission"`
	Times   []x4GameInfoTimeRuleData `xml:"data"`
}

type x4GameInfoTimeRuleData struct {
	XMLName xml.Name `xml:"data"`
	Minutes int      `xml:"time,attr"`
}
