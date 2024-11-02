package test

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"ssemu/internal/app"
	auth "ssemu/internal/app/auth/handlers"
	game "ssemu/internal/app/game/handlers"
	"ssemu/internal/domain"
	"ssemu/internal/network"
	"testing"
	"time"
)

type TestClient struct {
	SessionId network.SessionId
	UserId    domain.UserId
	Username  string
	Logger    *slog.Logger

	tester           *testing.T
	conn             net.Conn
	packetBufferList *list.List
	encrypt          bool
}

func NewTestClient(t *testing.T, protocol string, port uint16) TestClient {
	host := os.Getenv("EMU_TESTING_HOST")
	if len(host) == 0 {
		host = "localhost"
	}
	logger, err := CreateLogger(os.Getenv("EMU_TESTING_LOGS"))
	if err != nil {
		t.Error("create logger failed", err)
		t.FailNow()
	}
	logger = logger.With(
		slog.String("test", t.Name()),
		slog.Int("port", int(port)),
	)
	hostEndpoint := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout(protocol, hostEndpoint, 5*time.Second)
	if err != nil {
		logger.Error("server connection timeout", "error", err)
		t.Error("server connection timeout")
		t.FailNow()
	}
	tc := TestClient{
		SessionId:        0,
		UserId:           0,
		Username:         "",
		Logger:           logger,
		tester:           t,
		conn:             conn,
		packetBufferList: list.New(),
		encrypt:          port == app.GameServerSettings.Port,
	}
	tc.packetBufferList.Init()
	return tc
}

func (t *TestClient) recoverPanic() {
	if r := recover(); r != nil {
		t.Logger.Error("client panic", "err", r)
		t.tester.Error("client panic")
		t.tester.FailNow()
	}
}

func (t *TestClient) Receive(id network.PacketId) *network.Packet {
	defer t.recoverPanic()

	if packet, err := t.fetchFromPacketBufferList(id); err == nil {
		return packet
	}

	buffer := make([]byte, 4096)

	t.Logger.Debug("receiving packets...")
	bytesRead, err := t.conn.Read(buffer)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			t.Logger.Error("server connection closed", "err", err, "bytesRead", bytesRead)
			t.tester.Error("server connection closed")
			t.tester.FailNow()
		}
		t.Logger.Error("failed to receive packet", "err", err, "bytesRead", bytesRead)
		t.tester.Error("failed to receive packet")
		t.tester.FailNow()
	}

	t.Logger.Debug("reading packet list...")
	r := bytes.NewReader(buffer)
	for {
		var packetSize uint16
		if err := binary.Read(r, binary.LittleEndian, &packetSize); err != nil {
			break
		}
		if packetSize == 0 {
			break
		}
		packetBuffer := make([]byte, 512)
		binary.LittleEndian.PutUint16(packetBuffer[:2], packetSize)
		if err := binary.Read(r, binary.LittleEndian, packetBuffer[2:packetSize]); err != nil {
			t.Logger.Error("binary read packet failed", "err", err.Error())
			break
		}
		packet := network.GetPacket(packetBuffer)
		t.packetBufferList.PushBack(&packet)
		t.Logger.Debug("packet received",
			"packet", packet.Id().HexString(),
			"packetSize", packetSize,
			"packetList", PacketListToString(t.packetBufferList),
		)
	}

	packet, err := t.fetchFromPacketBufferList(id)
	if err != nil {
		t.Logger.Error("packet not found in packet list", "packet", id.HexString())
		t.tester.Error("packet not found in packet list")
		t.tester.FailNow()
	}

	return packet
}

func (t *TestClient) fetchFromPacketBufferList(id network.PacketId) (*network.Packet, error) {
	t.Logger.Debug("searching packet list...")
	for e := t.packetBufferList.Front(); e != nil; e = e.Next() {
		packet := e.Value.(*network.Packet)
		if packet.Id() == id {
			t.packetBufferList.Remove(e)
			t.Logger.Debug("packet read from buffer",
				"packet", packet.Id().HexString(),
				"packetList", PacketListToString(t.packetBufferList),
			)
			return packet, nil
		}
	}
	return nil, errors.New("packet not found")
}

func (t *TestClient) Send(packet network.Packet) {
	buffer := packet.Data()
	bufferLen := len(buffer)
	if t.encrypt {
		network.Encrypt(buffer[4:bufferLen])
	}
	if _, err := t.conn.Write(buffer); err != nil {
		t.Logger.Error("failed to send packet", "error", err)
		t.tester.Error("failed to send packet")
		t.tester.FailNow()
		return
	}
	t.Logger.Debug("packet sent", "packet", packet.Id().HexString(), "length", bufferLen, "encrypt", t.encrypt)
}

func (t *TestClient) Dispose() {
	t.conn.Close()
	t.Logger.Debug("client disposed")
}

func (t *TestClient) LoginAuth(username string, password string) auth.LoginResult {
	t.Logger = t.Logger.With("username", username)
	t.Logger.Info("auth login...", "password", password)
	req := network.NewPacket(auth.LoginReq)
	defer req.Free()
	req.WriteStringSlice(username, 13)
	req.WriteStringSlice(password, 13)
	t.Send(req)
	ack := t.Receive(auth.LoginAck)
	t.SessionId = network.SessionId(ack.ReadU32())
	ack.Skip(12)
	result := auth.LoginResult(ack.ReadU8())
	t.Logger.Debug("auth ack",
		"result", result,
		"session", t.SessionId,
	)
	t.Username = username
	return result
}

func (t *TestClient) LoginGame(username string, authSessionId network.SessionId) game.LoginResult {
	t.Logger = t.Logger.With("username", username)
	t.Logger.Info("game login...", "authSessionId", authSessionId)
	req := network.NewPacket(game.LoginReq)
	defer req.Free()
	req.WriteStringSlice(username, 13)
	req.Skip(30)
	req.WriteU32(uint32(authSessionId))
	t.Send(req)
	ack := t.Receive(game.LoginAck)
	t.UserId = domain.UserId(ack.ReadU32())
	ack.Skip(4)
	result := game.LoginResult(ack.ReadU8())
	t.Logger = t.Logger.With("user", t.UserId)
	t.Logger.Debug("game ack", "result", result)
	if result == game.Ok {
		t.Receive(game.LicenseDataAck)
		ack = t.Receive(game.CharacterSlotDataAck)
		charCount := ack.ReadU8()
		for i := 0; i < int(charCount); i++ {
			t.Receive(game.CharacterDataAck)
			t.Receive(game.CharacterEquipDataAck)
		}
		t.Receive(game.InventoryDataAck)
		t.Receive(game.ServerResultAck)
		t.Receive(game.StatsDataAck)
		t.Receive(game.ServerResultAck)
		t.Logger.Debug("player data received")
	}
	return result
}
