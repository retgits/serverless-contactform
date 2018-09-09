// Package util implements utility methods
package util

// The imports
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// GetSSMParameter gets a parameter from the AWS Simple Systems Manager service.
func GetSSMParameter(awsSession *session.Session, name string, decrypt bool) (string, error) {
	// Create an instance of the SSM Session
	ssmSession := ssm.New(awsSession)

	// Create the request to SSM
	getParameterInput := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(decrypt),
	}

	// Get the parameter from SSM
	param, err := ssmSession.GetParameter(getParameterInput)
	if err != nil {
		return "", err
	}

	return *param.Parameter.Value, nil
}
