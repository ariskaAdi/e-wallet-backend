package auth

import (
	"ariskaAdi/e-wallet/infra/response"
	"ariskaAdi/e-wallet/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthEntity struct {
	Id        int    `db:"id"`
	Username  string `db:"username"`
	Email     string `db:"email"`
	PublicId  uuid.UUID `db:"public_id"`
	Password  string `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewFormRegisterRequset(req RegisterRequestPayload) AuthEntity {
	return AuthEntity{
		Username: req.Username,
		Email: req.Email,
		PublicId: uuid.New(),
		Password: req.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewFormLoginRequset(req LoginRequestPayload) AuthEntity {
	return AuthEntity{
		Email: req.Email,
		Password: req.Password,
	}
}

func (a AuthEntity) Validate() (err error) {
	if err = a.EmailValidate();err != nil {
		return	
	}

	if err = a.PasswordValidate();err != nil {
		return
	}

	return
}

func (a AuthEntity) EmailValidate() (err error) {
	if a.Email == "" {
		return response.ErrEmailRequired
	} 

	emails := strings.Split(a.Email, "@")
	if len(emails) != 2 {
		return response.ErrEmailInvalid
	}

	return
}

func (a AuthEntity) PasswordValidate() (err error) {
	if a.Password == "" {
		return response.ErrPasswordRequired
	}

	if len(a.Password)< 6 {
		return response.ErrPasswordInvalid
	}
	return
}

func (a *AuthEntity) EncryptPassword(sal int) (err error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	a.Password = string(encryptedPassword)
	return nil
}

func (a AuthEntity) IsExist() bool {
	return  a.Id != 0 
}

func (a AuthEntity) VerifyPasswordFromEncrypted(plain string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(plain))
}

func (a AuthEntity) VerifyPasswordFromPlain(encrypted string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(a.Password))
}

func (a AuthEntity) GenerateToken(secret string) (tokenString string, err error) {
	return utils.GenerateToken(a.PublicId.String(),secret)
}