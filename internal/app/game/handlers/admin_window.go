package handlers

import (
	"context"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AdminShowWindowReqHandler struct{}

const AdminShowWindowReq network.PacketId = 0x46
const AdminShowWindowAck network.PacketId = 0x47

type AdminShowWindowResult byte

const (
	Allowed    AdminShowWindowResult = 0
	NotAllowed AdminShowWindowResult = 1
)

func (h *AdminShowWindowReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	ack := network.NewPacket(AdminShowWindowAck)
	defer ack.Free()

	isAdmin, err := dbIsUserAdmin(ctx, s.UserId())
	if err != nil {
		return err
	}
	if isAdmin {
		ack.WriteU8(byte(Allowed))
	} else {
		ack.WriteU8(byte(NotAllowed))
	}

	span.SetAttributes(attribute.Key("isAdmin").Bool(isAdmin))

	s.Send(&ack)

	return nil
}

func dbIsUserAdmin(ctx context.Context, id domain.UserId) (bool, error) {
	ctx, span := database.GetTracer().Start(ctx, "dbIsUserAdmin")
	defer span.End()
	const query string = "SELECT is_admin FROM users u WHERE u.id = ? LIMIT 1"
	db := database.GetConn()
	r := db.QueryRowContext(ctx, query, id)
	var isAdmin bool
	if err := r.Scan(&isAdmin); err != nil {
		return false, err
	}
	return isAdmin, nil
}

func (h *AdminShowWindowReqHandler) AllowAnonymous() bool { return false }

func (h *AdminShowWindowReqHandler) GetHandlerName() string { return "AdminShowWindowReqHandler" }
