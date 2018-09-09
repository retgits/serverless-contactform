// Package util implements utility methods
package util

// The imports
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SendEmail sends an email using Simple Email Service.
func SendEmail(awsSession *session.Session, toAddress string, bodyContent string, fromAddress string, subject string) error {
	// Create an instance of the SES Session
	sesSession := ses.New(awsSession)

	// Create the Email request
	sesEmailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(toAddress)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(bodyContent),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(fromAddress),
		ReplyToAddresses: []*string{
			aws.String(fromAddress),
		},
	}

	// Send the email
	_, err := sesSession.SendEmail(sesEmailInput)
	return err
}
