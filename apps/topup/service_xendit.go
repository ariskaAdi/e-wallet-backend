package topup

import (
	"ariskaAdi/e-wallet/internal/config"
	"context"
	"errors"

	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
)

type XenditService interface {
	GenerateSnapURL(ctx context.Context, model *TopUpEntity) error
	CheckPayment(ctx context.Context, invoiceId string) (string,error)
}

type xenditService struct {
	client *xendit.APIClient
}

func newXenditService(cfg *config.Config) XenditService {
	client := xendit.NewClient(cfg.Xendit.Key)

	return &xenditService{
		client: client,
	}
}

func (x *xenditService) GenerateSnapURL(ctx context.Context, model *TopUpEntity) error {

	var desc = "Topup e-wallet"
	var Currency = "IDR"
	req := invoice.CreateInvoiceRequest{
		ExternalId: model.TopUpId,
		Amount: float64(model.Amount),
		Description: &desc,
		Currency: &Currency,
	}

	
	resp, _, err := x.client.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(req).
		Execute()
	if err != nil {
		return err
	}

	model.AttachSnapUrl(resp.InvoiceUrl)
	return nil
}

func (x *xenditService) CheckPayment(
	ctx context.Context,
	invoiceId string,
) (string, error) {

	if invoiceId == "" {
		return  "", errors.New("invoice id is empty")
	}

	resp, _, err := x.client.InvoiceApi.GetInvoiceById(ctx, invoiceId).Execute()
	if err != nil {
		return  "", err
	}

	if resp.Status != invoice.INVOICESTATUS_EXPIRED {
        return string(resp.Status), nil
    }

    return "UNKNOWN", nil
}


