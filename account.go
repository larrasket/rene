package main

import (
	"errors"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var client *twitter.Client

func AddAccount(username string) (err error) {

	requestToken, err := login()
	if err != nil {
		logger.Error("Request Token Phase", err)
		return
	}

	tok, sec, err := receivePIN(requestToken)
	if err != nil {
		logger.Error("Access Token Phase", err)
		return
	}
	token := oauth1.NewToken(tok, sec)
	httpClient := AuthConfig.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)
	u, res, err := client.Accounts.
		VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		var p []byte
		res.Body.Read(p)
		logger.Error(`Something wrong happened in verfiny credentails`, err,
			res.Request.URL.String(), res.StatusCode, string(p))
		return
	}
	if u.ScreenName != username {
		err = errors.New(`username doesn't match`)
		logger.Error(fmt.Sprintf(`The provided username %s doesn't match %s`,
			username, u.ScreenName), err)
	}
	return
}

func login() (requestToken string, err error) {
	requestToken, _, err = AuthConfig.RequestToken()
	if err != nil {
		return "", err
	}

	authorizationURL, err := AuthConfig.AuthorizationURL(requestToken)
	if err != nil {
		return "", err
	}

	fmt.Printf("Open this URL in your browser:\n%s\n",
		authorizationURL.String())
	return requestToken, err
}

func receivePIN(requestToken string) (tok string, sec string, err error) {
	fmt.Printf("Paste your PIN here: ")
	var verifier string
	fmt.Scanf("%s", &verifier)
	tok, sec, err = AuthConfig.AccessToken(requestToken,
		"secret does not matter", verifier)
	return
}
