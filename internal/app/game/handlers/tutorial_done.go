package handlers

import (
	"context"
	"ssemu/internal/database"
	"ssemu/internal/domain"
	"ssemu/internal/network"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TutorialDoneReqHandler struct{}

const TutorialDoneReq network.PacketId = 0x54

func (h *TutorialDoneReqHandler) HandlePacket(ctx context.Context, s network.Session, p network.Packet) error {
	span := trace.SpanFromContext(ctx)

	player, err := dbGetPlayer(ctx, s.UserId())
	if err != nil {
		return err
	}

	if player.TutorialStatus == domain.TutorialDone {
		return nil
	}

	player.TutorialStatus = domain.TutorialDone

	if err := dbUpdateTutorialStatus(ctx, player); err != nil {
		return err
	}

	span.SetAttributes(attribute.Bool("tutorialDone", true))

	return nil
}

func dbUpdateTutorialStatus(ctx context.Context, player domain.Player) error {
	ctx, span := database.GetTracer().Start(ctx, "dbUpdateTutorialStatus")
	defer span.End()
	const query string = "UPDATE players SET tutorial_status = ? WHERE user_id = ?"
	if _, err := database.GetConn().ExecContext(ctx, query, player.TutorialStatus, player.UserId); err != nil {
		return err
	}
	return nil
}

func (h *TutorialDoneReqHandler) AllowAnonymous() bool { return false }

func (h *TutorialDoneReqHandler) GetHandlerName() string { return "TutorialDoneReqHandler" }
