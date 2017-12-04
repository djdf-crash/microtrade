package utils

import (
	"net/smtp"
)

func SendEmail(serveSMTP, portServe, sender, passwordSender, recipient, body string) error {

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		sender,
		passwordSender,
		serveSMTP,
	)

	msg := "From: " + sender + "\n" +
		"To: " + recipient + "\n" +
		"Subject: Reset password\n\n" +
		body

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		serveSMTP+portServe,
		auth,
		sender,
		[]string{recipient},
		[]byte(msg),
	)
	if err != nil {
		return err
	}

	return nil
}
