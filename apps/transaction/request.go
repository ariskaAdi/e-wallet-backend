package transaction

type TransferInquiryRequestPayload struct {
	Dof string `json:"dof_number"`
}

type TransferExecuteRequest struct {
	InquiryKey  string `json:"inquiry_key"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}
