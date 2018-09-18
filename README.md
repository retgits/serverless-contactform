# contactform-lambda - A serverless app to serve the contactform

A serverless tool designed to serve a contactform.

## Layout
```bash
.                    
├── test            
│   └── event.json      <-- Sample event to test using SAM local
├── .gitignore          <-- Ignoring the things you do not want in git
├── LICENSE             <-- The license file
├── main.go             <-- The Lambda code
├── Makefile            <-- Makefile to build and deploy
├── README.md           <-- This file
└── template.yaml       <-- SAM Template
```

## Installing
There are a few ways to install this project

### Get the sources
You can get the sources for this project by simply running
```bash
$ go get -u github.com/retgits/contactform-lambda/...
```

### Deploy
Deploy the Lambda app by running
```bash
$ make deploy
```

## Parameters
### AWS Systems Manager parameters
The code will automatically retrieve the below list of parameters from the AWS Systems Manager Parameter store:

* **/google/recaptcha/secret**: Your reCAPTCHA Secret Token (check the Google reCAPTCHA documentation on how to get this parameter)
* **/google/recaptcha/email**: A validated email address in SES (check the Google reCAPTCHA documentation on how to get this parameter)

### Deployment parameters
In the `template.yaml` there are certain deployment parameters:

* **region**: The AWS region in which the code is deployed

## Make targets
contactform-lambda has a _Makefile_ that can be used for most of the operations

```
usage: make [target]
```

* **deps**: Gets all dependencies for this app
* **clean** : Removes the dist directory
* **build**: Builds an executable to be deployed to AWS Lambda
* **test-lambda**: Clean, builds and tests the code by using the AWS SAM CLI
* **deploy**: Cleans, builds and deploys the code to AWS Lambda