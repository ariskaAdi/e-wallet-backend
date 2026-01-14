package transaction

type TransferInquiryResponse struct {
	InquiryKey string `json:"inquiry_key"`
	Dof        string `json:"dof_number"`
	ExpiredAt  string `json:"expired_at"`
}

type TransferExecuteResponse struct {
	Amount int64 `json:"amount"`
}