package mail

import (
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     int
	User string
	Pass string
	From string
}

func SendOtp(
	cfg SMTPConfig,
	to string,
	username string,
	otp string,
) (err error) {
	auth := smtp.PlainAuth(
		"",
		cfg.User,
		cfg.Pass,
		cfg.Host,
	)

	subject := "Subject: Verifikasi Akun Anda\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(`
		<h3>Halo %s</h3>
		<p>Kode OTP kamu:</p>
		<h2>%s</h2>
		<p>Berlaku selama 5 menit.</p>
	`, username, otp)

	msg := []byte(subject + mime + body)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	return smtp.SendMail(
		addr,
		auth,
		cfg.From,
		[]string{to},
		msg,
	)
}