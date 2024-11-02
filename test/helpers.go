package test

import (
	"container/list"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"ssemu/internal/network"
	"ssemu/internal/resources"
	"strings"
	"testing"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

var resourcesLoaded bool

func LoadResources(t *testing.T) {
	if !resourcesLoaded {
		clientDirPath := os.Getenv("EMU_TESTING_CLIENT")
		if err := resources.Load(clientDirPath); err != nil {
			t.Errorf("Failed to load client resources. %v", err)
			t.FailNow()
			return
		}
		resourcesLoaded = true
	}
}

var logger *slog.Logger

func CreateLogger(logPath string) (*slog.Logger, error) {
	if logger != nil {
		return logger, nil
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o755)
	if err != nil {
		return nil, err
	}
	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger, nil
}

func PacketListToString(list *list.List) string {
	var sb strings.Builder
	for e := list.Front(); e != nil; e = e.Next() {
		packet := e.Value.(*network.Packet)
		sb.WriteString(fmt.Sprintf("%s ", packet.Id().HexString()))
	}
	return sb.String()
}
