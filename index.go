package main

// The imports
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
)

// Constants
const (
	// The URL to validate reCAPTCHA
	recaptchaURL = "https://www.google.com/recaptcha/api/siteverify"
)

// Variables
var (
	// The reCAPTCHA Secret Token
	recaptchaSecret = os.Getenv("RECAPTCHA_SECRET")
	// The email address to send data to
	emailAddress = os.Getenv("EMAIL_ADDRESS")
	// The email password to use
	emailPassword = os.Getenv("EMAIL_PASSWORD")
	// The SMTP server
	smtpServer = os.Getenv("SMTP_SERVER")
	// The SMTP server port
	smtpPort = os.Getenv("SMTP_PORT")
)

// Handler is the main entry point into tjhe function code as mandated by ZEIT
func Handler(w http.ResponseWriter, r *http.Request) {
	// HTTPS will do a PreFlight CORS using the OPTIONS method.
	// To complete that a special response should be sent
	if r.Method == http.MethodOptions {
		response(w, true, "", r.Method)
		return
	}

	// Parse the request body to a map
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	u, err := url.ParseQuery(buf.String())
	if err != nil {
		response(w, false, fmt.Sprintf("There was an error sending your form data: %s", err.Error()), r.Method)
		return
	}

	// Prepare the POST parameters
	urlData := url.Values{}
	urlData.Set("secret", recaptchaSecret)
	urlData.Set("response", u["g-recaptcha-response"][0])

	// Validate the reCAPTCHA
	resp, err := httpcall(recaptchaURL, "POST", "application/x-www-form-urlencoded", urlData.Encode(), nil)
	if err != nil {
		response(w, false, fmt.Sprintf("There was an error sending your form data: %s", err.Error()), r.Method)
		return
	}

	// Validate if the reCAPTCHA was successful
	if !resp.Body["success"].(bool) {
		response(w, false, fmt.Sprintf("There was an error sending your form data: %s", fmt.Sprintf("%v", resp.Body["error-codes"])), r.Method)
		return
	}

	// Set up email authentication information.
	auth := smtp.PlainAuth(
		"",
		emailAddress,
		emailPassword,
		smtpServer,
	)

	// Prepare the email
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := fmt.Sprintf("Subject: [BLOG] Message from %s %s!\n", u["name"][0], u["surname"][0])
	msg := []byte(fmt.Sprintf("%s%s\n%s\n\n%s", subject, mime, u["message"][0], u["email"][0]))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err = smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpServer, smtpPort),
		auth,
		emailAddress,
		[]string{emailAddress},
		msg,
	)
	if err != nil {
		fmt.Printf("[BLOG] Message from %s %s\n%s\n%s\nThe message was not sent: %s", u["name"][0], u["surname"][0], u["message"][0], u["email"][0], err.Error())
		response(w, false, "There was an error sending your email, but we've logged the data...", r.Method)
		return
	}

	// Return okay response
	response(w, true, "Thank you for your email! I'll contact you soon.", r.Method)
	return
}

func response(w http.ResponseWriter, success bool, message string, method string) {
	// Create a map for the response body
	body := make(map[string]interface{})

	// Prepare the return data
	if success {
		body["type"] = "success"
	} else {
		body["type"] = "danger"
	}
	body["message"] = message
	bodyString, _ := json.Marshal(body)

	// Return the response
	if method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "*")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyString)
}

// HTTPResponse is the response type for the HTTP requests
type HTTPResponse struct {
	Body       map[string]interface{}
	StatusCode int
	Headers    http.Header
}

// httpcall executes an HTTP request request to a URL and returns the response body as a JSON object
func httpcall(URL string, requestType string, encoding string, payload string, header http.Header) (HTTPResponse, error) {
	// Instantiate a response object
	httpresponse := HTTPResponse{}

	// Prepare placeholders for the request and the error object
	req := &http.Request{}
	var err error

	// Create a request
	if len(payload) > 0 {
		req, err = http.NewRequest(requestType, URL, strings.NewReader(payload))
		if err != nil {
			return httpresponse, fmt.Errorf("error while creating HTTP request: %s", err.Error())
		}
	} else {
		req, err = http.NewRequest(requestType, URL, nil)
		if err != nil {
			return httpresponse, fmt.Errorf("error while creating HTTP request: %s", err.Error())
		}
	}

	// Associate the headers with the request
	if header != nil {
		req.Header = header
	}

	// Set the encoding
	if len(encoding) > 0 {
		req.Header["Content-Type"] = []string{encoding}
	}

	// Execute the HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return httpresponse, fmt.Errorf("error while performing HTTP request: %s", err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return httpresponse, err
	}

	httpresponse.Headers = res.Header
	httpresponse.StatusCode = res.StatusCode

	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		return httpresponse, fmt.Errorf("error while unmarshaling HTTP response to JSON: %s", err.Error())
	}

	httpresponse.Body = data

	return httpresponse, nil
}
