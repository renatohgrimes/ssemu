package channels

import (
	"fmt"
	"ssemu/internal/domain"
	"ssemu/internal/resources"
)

var channels []Channel

func Load() {
	if channels != nil {
		return
	}
	for _, x4Channel := range resources.GetChannels() {
		// ignore non free channels since we do not support player levels
		if x4Channel.Type != 10 {
			continue
		}
		channel := NewChannel(uint32(x4Channel.Id), x4Channel.Language.Nations[0].Code)
		channels = append(channels, channel)
	}
}

func List() []Channel { return channels }

func GetChannelById(channelId uint32) (*Channel, error) {
	for _, channel := range channels {
		if channel.id == channelId {
			return &channel, nil
		}
	}
	return nil, fmt.Errorf("server channel %d not found", channelId)
}

func GetChannelByName(name string) (*Channel, error) {
	for _, channel := range channels {
		if channel.name == name {
			return &channel, nil
		}
	}
	return nil, fmt.Errorf("server channel %s not found", name)
}

func RemoveUser(userId domain.UserId) {
	for _, channel := range channels {
		channel.Leave(userId)
	}
}

func GetUserChannel(userId domain.UserId) (*Channel, error) {
	for _, channel := range channels {
		if _, exists := channel.sessions[userId]; exists {
			return &channel, nil
		}
	}
	return nil, fmt.Errorf("user %d channel not found", userId)
}
