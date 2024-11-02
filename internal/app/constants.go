package app

import (
	"ssemu/internal/database"
	"ssemu/internal/network"
)

var Version = "dev"

const MaxCapacity int = 1024

var AuthServerSettings = network.ServerSettings{
	Name:        "ssemu.auth",
	Port:        28002,
	Type:        network.Auth,
	MaxCapacity: MaxCapacity,
	Version:     Version,
}

var GameServerSettings = network.ServerSettings{
	Name:        "ssemu.game",
	Port:        28008,
	Type:        network.Game,
	MaxCapacity: MaxCapacity,
	Version:     Version,
}

var ChatServerSettings = network.ServerSettings{
	Name:        "ssemu.chat",
	Port:        28012,
	Type:        network.Chat,
	MaxCapacity: MaxCapacity,
	Version:     Version,
}

var RelayServerSettings = network.ServerSettings{
	Name:        "ssemu.relay",
	Port:        28013,
	Type:        network.Relay,
	MaxCapacity: MaxCapacity,
	Version:     Version,
}

var Nat1ServerSettings = network.ServerSettings{
	Name:        "ssemu.nat1",
	Port:        38912,
	Type:        network.Nat,
	MaxCapacity: MaxCapacity,
	Version:     Version,
}

var DbSettings = database.Settings{
	Dsn:          "file:database.sqlite3?_journal=WAL&_timeout=5000",
	MaxOpenConns: 4,
}
