package notificationservice

import (
	"argus-backend/internal/logger"
	"gopkg.in/gomail.v2"
)

func (wns *NotificationService) SendEmailNotification(text string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "just.sparkless@gmail.com")
	mail.SetHeader("To", "just.sparky@mail.ru")
	mail.SetHeader("Subject", "Argus notification")
	mail.SetBody("text/plain", text)

	d := gomail.NewDialer("smtp.gmail.com", 587, "just.sparkless@gmail.com", "")
	if err := d.DialAndSend(mail); err != nil {
		logger.Error("error sending notification via email: " + err.Error())
		return err
	}

	return nil
}
