package resources

import (
	"os"
	"path"
	"ssemu/internal/domain"
)

var cypherKey []byte
var x4 x4Manager

func Load(clientDirPath string) error {
	var err error
	cypherKey, err = extractCypherKey(clientDirPath)
	if err != nil {
		return err
	}
	hdPath := path.Join(clientDirPath, "resource.s4hd")
	resourcesDirPath := path.Join(clientDirPath, "_resources")
	files, err := extractX4Files(hdPath, resourcesDirPath, cypherKey)
	if err != nil {
		return err
	}
	x4, err = unmarshalX4(files)
	if err != nil {
		return err
	}
	return nil
}

func extractCypherKey(clientDirPath string) ([]byte, error) {
	exePath := path.Join(clientDirPath, "S4Client.exe")
	file, err := os.Open(exePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if _, err = file.Seek(0x51FC00, 0); err != nil {
		return nil, err
	}
	data := make([]byte, 800)
	if _, err = file.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}

func IsMapValid(mapId int) bool {
	for _, mapValue := range x4.gameInfo.Maps.Values {
		if mapValue.Id == mapId {
			return true
		}
	}
	return false
}

func IsGameLimitValid(score int, timeMinutes int, scores []x4GameInfoScoreRuleData, times []x4GameInfoTimeRuleData) bool {
	scoreIndex := 0
	timeIndex := 0
	for i, limit := range scores {
		if limit.Score == int(score) {
			scoreIndex = i
		}
	}
	for i, limit := range times {
		if limit.Minutes == int(timeMinutes) {
			timeIndex = i
		}
	}
	if scoreIndex == 0 || timeIndex == 0 {
		return false
	}
	if scoreIndex == timeIndex {
		return true
	}
	return false
}

func GetDefaultItems(gender domain.CharacterGender) []x4DefaultItem {
	if gender == domain.Male {
		return x4.defaultItems.Male.Items
	}
	return x4.defaultItems.Female.Items
}

func GetCypherKey() []byte { return cypherKey }

func GetChannelCapacity() int { return x4.channels.Setting.ChannelCapacity }

func GetChannels() []x4ChannelInfo { return x4.channels.Channels }

func GetGameScoreLimit() x4GameInfoScore { return x4.gameInfo.Scores }

func GetGameTimeLimit() x4GameInfoTime { return x4.gameInfo.Times }
