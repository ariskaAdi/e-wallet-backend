package topup

type CreateTopUpResponse struct {
	TopUpId string `json:"topup_id"`
	SnapURL string `json:"snap_url"`
	Status  string `json:"status"`
}
