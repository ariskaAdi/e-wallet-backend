package auth

import (
	infrafiber "ariskaAdi/e-wallet/infra/fiber"
	"ariskaAdi/e-wallet/infra/response"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	svc service
}

func newHandler(svc service) handler {
	return handler{svc: svc}
}

func (h handler) register(ctx *fiber.Ctx) error {
	var req = RegisterRequestPayload{}

	if err := ctx.BodyParser(&req); err != nil {
		myErr := response.ErrorBadRequest
		return infrafiber.NewResponse(
			infrafiber.WithMessage(err.Error()),
			infrafiber.WithError(myErr),
			infrafiber.WithHttpCode(http.StatusBadRequest),
			infrafiber.WithMessage("register failed"),
		).Send(ctx)
	}

	if err := h.svc.register(ctx.UserContext(), req); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		return infrafiber.NewResponse(
			infrafiber.WithMessage(err.Error()),
			infrafiber.WithError(myErr),
		).Send(ctx)
	}

	return infrafiber.NewResponse(
		infrafiber.WithHttpCode(http.StatusCreated),
		infrafiber.WithMessage("register success"),
	).Send(ctx)
	
}

func (h handler) login(ctx *fiber.Ctx) error {
	var req = LoginRequestPayload{}

	if err := ctx.BodyParser(&req); err != nil {
		myErr := response.ErrorBadRequest
		return infrafiber.NewResponse(
			infrafiber.WithMessage(err.Error()),
			infrafiber.WithError(myErr),
			infrafiber.WithMessage("login failed"),
		).Send(ctx)
	}

	token, err := h.svc.login(ctx.UserContext(), req)
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

	return infrafiber.NewResponse(
		infrafiber.WithHttpCode(http.StatusCreated),
		infrafiber.WithPayload(map[string]interface{}{
			"access_token": token,
		}),
		infrafiber.WithMessage("login success"),
	).Send(ctx)
}

func (h handler) verifyOtp(ctx *fiber.Ctx) error {
	var req = ValidateOtpRequestPayload{}

	if err := ctx.BodyParser(&req); err != nil {
		myErr := response.ErrorBadRequest
		return infrafiber.NewResponse(
			infrafiber.WithMessage(err.Error()),
			infrafiber.WithError(myErr),
		).Send(ctx)
	}

	if err := h.svc.verifyOtp(ctx.UserContext(), req); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		return infrafiber.NewResponse(
			infrafiber.WithMessage(err.Error()),
			infrafiber.WithError(myErr),
		).Send(ctx)
	}

	return infrafiber.NewResponse(
		infrafiber.WithHttpCode(http.StatusCreated),
		infrafiber.WithMessage("verify otp success"),
	).Send(ctx)
	
	
}