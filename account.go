package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var client *twitter.Client

type authorizer struct{}

func (a *authorizer) Add(req *http.Request) {}

func AddAccount(username string) error {
	requestToken, err := login()
	if err != nil {
		return fmt.Errorf("Request Token Phase: %s", err)
	}

	tok, sec, err := receivePIN(requestToken)
	if err != nil {
		return fmt.Errorf("Access Token Phase: %s", err)
	}
	token := oauth1.NewToken(tok, sec)

	httpClient := AuthConfig.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)

	u, res, err := client.Accounts.
		VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		resStr, _ := io.ReadAll(res.Body)

		return fmt.Errorf(
			`Something wrong happened in verfiny credentails: %s, response body:
%s status code: %s: `, err, string(resStr), res.StatusCode)
	}
	defer res.Body.Close()

	if u.ScreenName != username {
		return fmt.Errorf(`The provided username %s doesn't match %s`,
			username, u.ScreenName)
	}

	stmt := `INSERT INTO accounts (username, account_token, account_secret)
	VALUES ("%s", "%s", "%s")`
	_, err = Db.Exec(fmt.Sprintf(stmt, username, tok, sec))
	return err
}

func login() (requestToken string, err error) {
	requestToken, _, err = AuthConfig.RequestToken()
	if err != nil {
		return "", fmt.Errorf(`Couldn't request token: `, err)
	}

	authorizationURL, err := AuthConfig.AuthorizationURL(requestToken)
	if err != nil {
		return "", fmt.Errorf(`Couldn't get authorization URL`, err)
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
