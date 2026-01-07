package auth

type RegisterRequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ValidateOtpRequestPayload struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}