package wallet

import (
	infrafiber "ariskaAdi/e-wallet/infra/fiber"
	"ariskaAdi/e-wallet/infra/response"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	svc service
}

func newHandler(svc service) handler {
	return handler{svc: svc}
}

func (h handler) GetMyWallet(ctx *fiber.Ctx) error {
	userPublicId := fmt.Sprintf("%v", ctx.Locals("user_public_id"))

	myWallet, err := h.svc.GetMyWallet(ctx.UserContext(), userPublicId)
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

	resp := MyWalletResponse{
		UserPublicId: myWallet.UserPublicId,
		WalletPublicId: myWallet.WalletPublicId,
		Name: myWallet.Name,
		Balance:      myWallet.Balance,
		CreatedAt:    myWallet.CreatedAt,
		UpdatedAt:    myWallet.UpdatedAt,
	}

	return infrafiber.NewResponse(
		infrafiber.WithHttpCode(http.StatusOK),
		infrafiber.WithPayload(resp),
	).Send(ctx)
}

func (h handler) GetSomeoneWallet(ctx *fiber.Ctx) error {

	req := GetWalletByIdRequestPayload{}

	if err := ctx.BodyParser(&req); err != nil {
		myErr := response.ErrorBadRequest
		return infrafiber.NewResponse(
			infrafiber.WithMessage(err.Error()),
			infrafiber.WithError(myErr),
		).Send(ctx)
	}

	wallet, err := h.svc.GetWalletId(ctx.UserContext(), req)
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

	resp := FindWalletResponse{
		Name: wallet.Name,
	}

	return infrafiber.NewResponse(
		infrafiber.WithHttpCode(http.StatusOK),
		infrafiber.WithPayload(resp),
	).Send(ctx)
}