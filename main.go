// Package main is the main implementation of the contactform serverless app.
package main

// The imports
import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/retgits/lambda-util"
)

// Constants
const (
	// The name of the reCAPTCHA Secret Token parameter in Amazon SSM
	tokenName = "/google/recaptcha/secret"
	// The name of email address parameter in Amazon SSM
	emailTokenName = "/google/recaptcha/email"
	// The URL to validate reCAPTCHA
	recaptchaURL = "https://www.google.com/recaptcha/api/siteverify"
)

// Variables
var (
	// The region in which the Lambda function is deployed
	awsRegion = util.GetEnvKey("region", "us-west-2")
)

// The handler function is executed every time that a new Lambda event is received.
// It takes a JSON payload (you can see an example in the event.json file) and only
// returns an error if the something went wrong. The event comes fom API Gateway
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Create a headers map to enable CORS
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"

	// Create a new session without AWS credentials.
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	// Get the reCAPTCHA Secret Token
	recaptchaSecret, err := util.GetSSMParameter(awsSession, tokenName, true)
	if err != nil {
		return errorResponse(headers, fmt.Sprintf("There was an error sending your form data: %s", err.Error())), nil
	}

	// Get the email address
	emailAddress, err := util.GetSSMParameter(awsSession, emailTokenName, true)
	if err != nil {
		return errorResponse(headers, fmt.Sprintf("There was an error sending your form data: %s", err.Error())), nil
	}

	// Parse the request body to a map
	u, err := url.ParseQuery(request.Body)
	if err != nil {
		return errorResponse(headers, fmt.Sprintf("There was an error sending your form data: %s", err.Error())), nil
	}

	// Prepare the POST parameters
	urlData := url.Values{}
	urlData.Set("secret", recaptchaSecret)
	urlData.Set("response", u["g-recaptcha-response"][0])

	// Validate the reCAPTCHA
	resp, err := util.HTTPPost(recaptchaURL, "application/x-www-form-urlencoded", urlData)
	if err != nil {
		return errorResponse(headers, fmt.Sprintf("There was an error sending your form data: %s", err.Error())), nil
	}

	// Validate if the reCAPTCHA was successful
	if !resp.Body["success"].(bool) {
		return errorResponse(headers, fmt.Sprintf("There was an error sending your form data: %s", fmt.Sprintf("%v", resp.Body["error-codes"]))), nil
	}

	// Send email
	err = util.SendEmail(awsSession, emailAddress, fmt.Sprintf("%s\n\n%s", u["message"][0], u["email"][0]), emailAddress, fmt.Sprintf("[BLOG] Message from %s %s", u["name"][0], u["surname"][0]))
	if err != nil {
		fmt.Printf("[BLOG] Message from %s %s\n%s\n%s\nThe message was not sent: %s", u["name"][0], u["surname"][0], u["message"][0], u["email"][0], err.Error())
		return errorResponse(headers, "There was an error sending your email, but we've logged the data..."), nil
	}

	// Return okay response
	return okayResponse(headers, "Thank you for your email! I'll contact you soon."), nil
}

// The main method is executed by AWS Lambda and points to the handler
func main() {
	lambda.Start(handler)
}

func errorResponse(headers map[string]string, reason string) events.APIGatewayProxyResponse {
	// Create a map for the response body
	body := make(map[string]interface{})

	// Prepare the return data
	body["type"] = "danger"
	body["message"] = reason
	bodyString, _ := json.Marshal(body)

	// Return the response
	return events.APIGatewayProxyResponse{Body: string(bodyString), StatusCode: 200, Headers: headers}
}

func okayResponse(headers map[string]string, reason string) events.APIGatewayProxyResponse {
	// Create a map for the response body
	body := make(map[string]interface{})

	// Prepare the return data
	body["type"] = "success"
	body["message"] = reason
	bodyString, _ := json.Marshal(body)

	// Return the response
	return events.APIGatewayProxyResponse{Body: string(bodyString), StatusCode: 200, Headers: headers}
}
