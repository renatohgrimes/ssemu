package resources

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type x4Manager struct {
	skills       x4ItemsRoot
	weapons      x4ItemsRoot
	costumes     x4ItemsRoot
	premiums     x4ItemsRoot
	shop         x4ItemsRoot
	channels     x4ChannelRoot
	gameInfo     x4GameInfoRoot
	defaultItems x4DefaultItemRoot
}

func unmarshalX4(files []x4File) (x4Manager, error) {
	var err error
	x4 = x4Manager{}
	for _, file := range files {
		switch file.name {
		case "s4_actions.x4":
			err = xml.Unmarshal(file.contents, &x4.skills)
		case "s4_channel_setting.x4":
			err = xml.Unmarshal(file.contents, &x4.channels)
		case "s4_charged_item.x4":
			err = xml.Unmarshal(file.contents, &x4.premiums)
		case "s4_costumes.x4":
			x4fix(&file)
			err = xml.Unmarshal(file.contents, &x4.costumes)
		case "s4_gameinfo.x4":
			err = xml.Unmarshal(file.contents, &x4.gameInfo)
		case "s4_shop.x4":
			x4fix(&file)
			err = xml.Unmarshal(file.contents, &x4.shop)
		case "s4_weapons.x4":
			err = xml.Unmarshal(file.contents, &x4.weapons)
		case "s4_default_item.x4":
			err = xml.Unmarshal(file.contents, &x4.defaultItems)
		}
		if err != nil {
			return x4, fmt.Errorf("failed to load %s xml. %w", file.name, err)
		}
	}
	return x4, nil
}

func x4fix(file *x4File) {
	str := string(file.contents)
	str = strings.ReplaceAll(str, "fasle", "false")
	str = strings.Replace(str, "EUC-KR", "UTF-8", 1)
	file.contents = []byte(str)
}
