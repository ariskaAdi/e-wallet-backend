package topup

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TopUpStatus int

const (
	TopUpPending TopUpStatus = 1
	TopUpSuccess TopUpStatus = 2
	TopUpFailed TopUpStatus = 3
)

const (
	TO_PENDING = "pending"
	TO_SUCCESS = "success"
	TO_FAILED  = "failed"
)

var MappingTopUpStatus = map[TopUpStatus]string{
	TopUpPending: TO_PENDING,
	TopUpSuccess: TO_SUCCESS,
	TopUpFailed:  TO_FAILED,
}

type TopUpEntity struct {
	Id        int       `db:"id"`
	TopUpId   string    `db:"topup_id"`
	UserPublicId string `db:"user_public_id"`
	Amount    int64     `db:"amount"`
	Status    TopUpStatus `db:"status"`
	SnapURL string `db:"snap_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewTopUp(amount int64, userPublicId string) TopUpEntity {
	return TopUpEntity{
		TopUpId:   uuid.NewString(),
		UserPublicId: userPublicId,
		Amount:    amount,
		Status:    TopUpPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *TopUpEntity) AttachSnapUrl(url string) {
	t.SnapURL = url
	t.UpdatedAt = time.Now()
}

func (t *TopUpEntity) MarkSuccess() (err error) {
	if t.Status != TopUpPending {
		return errors.New("Topup is not pending")
	}

	t.Status = TopUpSuccess
	t.UpdatedAt = time.Now()
	return
}

func (t *TopUpEntity) MarkFailed() {
	t.Status = TopUpFailed
	t.UpdatedAt = time.Now()
}

func (s TopUpStatus) String() string {
	if v, ok := MappingTopUpStatus[s]; ok {
		return v
	}
	return "unknown"
}
