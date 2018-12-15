# serverless-contactform - A serverless app to serve the contactform

A serverless app designed to serve a contactform, running on [ZEIT](https://zeit.co)

## Layout

```bash
.
├── .env_template <-- A template file with the environment variables needed for the function
├── .gitignore    <-- Ignoring the things you do not want in git
├── go.mod
├── go.sum
├── index_test.go <-- Tests for the function code
├── index.go      <-- The actual function code
├── LICENSE       <-- The license file
├── Makefile      <-- Makefile to build and deploy
├── now.json      <-- Deployment descriptor for ZEIT
└── README.md     <-- This file :)
```

## Installing

### Get the sources

You can get the sources for this project by simply running

```bash
go get -u github.com/retgits/serverless-contactform/...
```

### Update secrets

Update the secrets by running

```bash
make update-secrets
```

This command _will_ delete and recreate the secrets. The secrets used are:

* RECAPTCHA_SECRET: The reCAPTCHA secret token you can get from the _Server side integration_ step in [Google reCAPTCHA](https://www.google.com/recaptcha)
* EMAIL_ADDRESS: The email address to send data to (like _you@example.com_)
* EMAIL_PASSWORD: The password needed to log in to the SMTP server
* SMTP_SERVER: The SMTP server
* SMTP_PORT: The SMTP server port

### Deploy

Deploy the app by running

```bash
make deploy
```