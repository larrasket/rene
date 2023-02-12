package main

import (
	"context"

	"github.com/dghubble/oauth1"
)

func ListenAndServe() {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	for k := range Accounts {
		token := oauth1.NewToken(Accounts[0].tok, Accounts[0].sec)
		Accounts[k].client = AuthConfig.Client(oauth1.NoContext, token)
		select {
		case <-ctx.Done():
			logger.Error(`Something wrong happened while reading from DMs`,
				context.Cause(ctx))
		default:
		}
		go DbDump(&Accounts[k], ctx, cancel) // add DMs to database
	}
}
