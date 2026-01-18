package topup

import (
	infrafiber "ariskaAdi/e-wallet/infra/fiber"
	"ariskaAdi/e-wallet/infra/response"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	svc Service
}

func newHandler(svc Service) handler {
	return handler{svc: svc}
}

func (h handler) CreateTopUp(ctx *fiber.Ctx) error {
	userPublicId := fmt.Sprintf("%v", ctx.Locals("user_public_id"))
	var req CreateTopUpRequest

	if err := ctx.BodyParser(&req) ; err != nil {
		myErr := response.ErrorBadRequest
		return infrafiber.NewResponse(
			infrafiber.WithMessage(err.Error()),
			infrafiber.WithError(myErr),
		).Send(ctx)
	}

	topupEntity, err := h.svc.CreateTopUp(ctx.UserContext(), req.Amount, userPublicId)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		return infrafiber.NewResponse(
			infrafiber.WithMessage(err.Error()),
			infrafiber.WithError(myErr),
		).Send(ctx)
	}

	resp := CreateTopUpResponse{
		TopUpId: topupEntity.TopUpId,
		SnapURL: topupEntity.SnapURL,
		Status:  topupEntity.Status.String(),
	}

	return infrafiber.NewResponse(
		infrafiber.WithHttpCode(http.StatusCreated),
		infrafiber.WithPayload(resp),
		infrafiber.WithMessage("create topup success"),
	).Send(ctx)
}

func (h handler) XenditCallback(ctx *fiber.Ctx) error {
	var payload map[string]interface{}

	if err := ctx.BodyParser(&payload); err != nil {
		myErr := response.ErrorBadRequest
		return infrafiber.NewResponse(
			infrafiber.WithMessage(err.Error()),
			infrafiber.WithError(myErr),
		).Send(ctx)
	}


	externalId, ok := payload["external_id"].(string)
	if !ok {
		myErr := response.ErrorBadRequest
		return infrafiber.NewResponse(
			infrafiber.WithMessage("external_id not found"),
			infrafiber.WithError(myErr),
		).Send(ctx)
	}

	status, _ := payload["status"].(string)

	switch status {
	case "PAID":
		_ = h.svc.ConfirmTopUp(ctx.Context(), externalId)
	case "EXPIRED":
		_ = h.svc.FailTopUp(ctx.Context(), externalId)
	}

	return ctx.SendStatus(http.StatusOK)
}