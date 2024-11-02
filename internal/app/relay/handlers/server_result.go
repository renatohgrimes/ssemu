package handlers

import (
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const ServerResultAck network.PacketId = 0x02

type ServerResult uint32

const (
	TunnelJoinSuccess  ServerResult = 3
	TunnelLeaveSuccess ServerResult = 6
)

func ackWriteServerResult(span trace.Span, ack *network.Packet, res ServerResult) {
	span.SetAttributes(attribute.Key("result").Int(int(res)))
	ack.WriteU32(uint32(res))
}
