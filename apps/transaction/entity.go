package transaction

import (
	"ariskaAdi/e-wallet/infra/response"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeCredit TransactionType = "CREDIT" 
	TransactionTypeDebit  TransactionType = "DEBIT"  
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusSuccess   TransactionStatus = "SUCCESS"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

type TransactionEntity struct {
	Id           		int `db:"id"`
	TransactionId 		string `db:"transaction_id"`
	WalletPublicId		string `db:"wallet_public_id"`
	SofNumber 			string `db:"sof_number"`
	DofNumber			string `db:"dof_number"`
	Type 				TransactionType `db:"type"`
	Amount 				int64 `db:"amount"`
	Status 				TransactionStatus `db:"status"`
	Reference 			string `db:"reference"`
	Description 		string `db:"description"`
	CreatedAt 			time.Time `db:"created_at"`
	UpdatedAt 			time.Time `db:"updated_at"`
}

func (t TransactionEntity) Validate() (err error){
	if t.Amount <= 0 {
		return response.ErrAmountInvalid
	}

	if t.Type != TransactionTypeCredit && t.Type != TransactionTypeDebit {
		return response.ErrTransactionTypeInvalid
	}

	return
}

func NewTransaction(walletPublicId  string) TransactionEntity {
	return TransactionEntity{
		WalletPublicId : walletPublicId ,
		Status:       TransactionStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func NewCreditTransaction(
	sourceWallet WalletEntity,
	destWallet WalletEntity,
	amount int64,
	description string,
) TransactionEntity {

	return TransactionEntity{
		TransactionId:  uuid.NewString(),
		WalletPublicId: destWallet.WalletPublicId,
		SofNumber:      sourceWallet.WalletPublicId,
		DofNumber:      destWallet.WalletPublicId,
		Type:           TransactionTypeCredit,
		Amount:         amount,
		Status:         TransactionStatusPending,
		Description:    description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}


func NewDebitTransaction(
	sourceWallet WalletEntity,
	destWallet WalletEntity,
	amount int64,
	description string,
) TransactionEntity {

	return TransactionEntity{
		TransactionId:  uuid.NewString(),
		WalletPublicId: sourceWallet.WalletPublicId,
		SofNumber:      sourceWallet.WalletPublicId,
		DofNumber:      destWallet.WalletPublicId,
		Type:           TransactionTypeDebit,
		Amount:         amount,
		Status:         TransactionStatusPending,
		Description:    description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
