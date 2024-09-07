package utils

import (
	"context"
	"lamhat/core"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func TriggerEmail(subject string, body string, recipient string) error {

	// Your available domain names can be found here:
	// (https://app.mailgun.com/app/domains)
	var yourDomain string = core.Config.EMAIL.MAILGUN_DOMAIN

	// You can find the Private API Key in your Account Menu, under "Settings":
	// (https://app.mailgun.com/app/account/security)
	var privateAPIKey string = core.Config.EMAIL.MAILGUN_API
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	// The message object allows you to add attachments and Bcc recipients
	var sender string = core.Config.EMAIL.MAINGUN_SENDER
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		core.Sugar.Error(err)
		return err
	}

	core.Sugar.Info("ID: %s Resp: %s\n", id, resp)
	return nil
}
