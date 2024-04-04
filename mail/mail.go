package mail

import (
	"bytes"
	"eth2-exporter/db"
	"eth2-exporter/templates"
	"eth2-exporter/types"
	"eth2-exporter/utils"
	"fmt"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailTemplate struct {
	Mail   types.Email
	Domain string
}

// SendMail sends an email to the given address with the given message.
// It will use smtp if configured otherwise it will use gunmail if configured.
func SendHTMLMail(to, subject string, msg types.Email, attachment []types.EmailAttachment) error {
	var renderer = templates.GetTemplate("mail/layout.html")

	var err error
	var body bytes.Buffer

	if utils.Config.Frontend.Mail.SendGrid.PrivateKey != "" {
		_ = renderer.ExecuteTemplate(&body, "layout", MailTemplate{Mail: msg, Domain: utils.Config.Frontend.SiteDomain})
		content := body.String()
		err = SendMailViaSendGrid(to, subject, content, createTextMessage(msg), attachment)
	} else {
		err = fmt.Errorf("Invalid config setting for mail-service")
	}
	return err
}

// SendMail sends an email to the given address with the given message.
// It will use smtp if configured otherwise it will use gunmail if configured.
func SendTextMail(to, subject, msg string, attachment []types.EmailAttachment) error {
	var err error
	if utils.Config.Frontend.Mail.SendGrid.PrivateKey != "" {
		err = SendMailViaSendGrid(to, subject, msg, "", attachment)
	} else {
		err = fmt.Errorf("invalid config for mail-service")
	}
	return err
}

func createTextMessage(msg types.Email) string {
	return fmt.Sprintf("%s\n\n%s\n\n― You are receiving this because you are staking on Ethermine Staking. You can manage your subscriptions at %s.", msg.Title, msg.Body, msg.SubscriptionManageURL)
}

// SendMailRateLimited sends an email to a given address with the given message.
// It will return a ratelimit-error if the configured ratelimit is exceeded.
func SendMailRateLimited(to, subject string, msg types.Email, attachment []types.EmailAttachment) error {
	if utils.Config.Frontend.MaxMailsPerEmailPerDay > 0 {
		now := time.Now()
		count, err := db.GetMailsSentCount(to, now)
		if err != nil {
			return err
		}
		if count >= utils.Config.Frontend.MaxMailsPerEmailPerDay {
			timeLeft := now.Add(utils.Day).Truncate(utils.Day).Sub(now)
			return &types.RateLimitError{TimeLeft: timeLeft}
		}
	}

	err := db.CountSentMail(to)
	if err != nil {
		// only log if counting did not work
		return fmt.Errorf("error counting sent email: %v", err)
	}

	err = SendHTMLMail(to, subject, msg, attachment)
	if err != nil {
		return err
	}

	return nil
}

// SendMailViaSendGrid to the given address with the given message using sendgrid
func SendMailViaSendGrid(toEmail, subject string, msgTxt string, msgHtml string, attachment []types.EmailAttachment) error {
	var err error
	from := mail.NewEmail(utils.Config.Frontend.SiteDomain, utils.Config.Frontend.Mail.SendGrid.NonReply)
	to := mail.NewEmail("Reciever Email", toEmail)
	message := mail.NewSingleEmail(from, subject, to, msgTxt, msgHtml)
	client := sendgrid.NewSendClient(utils.Config.Frontend.Mail.SendGrid.PrivateKey)
	response, err := client.Send(message)
	fmt.Println("Success in sending mail using sendgrid. Statuscode %w", response.StatusCode)
	return err
}
