package transaction

import (
	"time"

	"github.com/google/uuid"
)

type InquiryEntity struct {
	InquiryKey string    `db:"inquiry_key"`
	Sof        string    `db:"sof_number"`
	Dof        string    `db:"dof_number"`
	CreatedAt  time.Time `db:"created_at"`
	ExpiredAt  time.Time `db:"expired_at"`
}

func NewInquiry (req TransferInquiryRequestPayload, UserPublicId string) InquiryEntity {
	return InquiryEntity{
		InquiryKey: uuid.NewString(),
		Sof: 		UserPublicId,
		Dof:        req.Dof,
		CreatedAt:  time.Now(),
		ExpiredAt:  time.Now().Add(time.Minute * 5),
	}
}