package transaction

type TransferInquiryRequestPayload struct {
	DestinationWalletPublicId string `json:"destination_wallet_public_id"`
}

type TransferExecuteRequest struct {
	InquiryKey  string `json:"inquiry_key"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}
