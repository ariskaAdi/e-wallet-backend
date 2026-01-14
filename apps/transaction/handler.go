package transaction

import (
	infrafiber "ariskaAdi/e-wallet/infra/fiber"
	"ariskaAdi/e-wallet/infra/response"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	svc service
}

func newHandler(svc service) handler {
	return handler{svc: svc}
}

func (h handler) TransferInquiry(ctx *fiber.Ctx) error {
	userPublicId := fmt.Sprintf("%v", ctx.Locals("user_public_id"))

	var req = TransferInquiryRequestPayload{}
	if err := ctx.BodyParser(&req); err != nil {
		myErr := response.ErrorBadRequest
		return 	infrafiber.NewResponse(
				infrafiber.WithMessage(err.Error()),
				infrafiber.WithError(myErr),
		).Send(ctx)
	}

	inquiry, err := h.svc.TransferInquiry(ctx.UserContext(), req, userPublicId)
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

	resp := TransferInquiryResponse{
		InquiryKey : inquiry.InquiryKey,
		Dof : inquiry.Dof,
		ExpiredAt: inquiry.ExpiredAt.Format(time.RFC3339),
	}

	return infrafiber.NewResponse(
		infrafiber.WithHttpCode(http.StatusOK),
		infrafiber.WithMessage("transfer inquiry success"),
		infrafiber.WithPayload(resp),
	).Send(ctx)
}

func (h handler) TransferExecute(ctx *fiber.Ctx) error {
	userPublicId := fmt.Sprintf("%v", ctx.Locals("user_public_id"))

	var req = TransferExecuteRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		myErr := response.ErrorBadRequest
		return 	infrafiber.NewResponse(
				infrafiber.WithMessage(err.Error()),
				infrafiber.WithError(myErr),
		).Send(ctx)
	}

	_, err := h.svc.TransferExecute(ctx.UserContext(), req, userPublicId)
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

	res := TransferExecuteResponse{
		Amount: req.Amount,
	}

	return infrafiber.NewResponse(
		infrafiber.WithHttpCode(http.StatusOK),
		infrafiber.WithMessage("transfer execute success"),
		infrafiber.WithPayload(res),
	).Send(ctx)
}