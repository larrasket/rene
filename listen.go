package main

import (
	"context"

	"github.com/dghubble/oauth1"
)

func ListenAndServe() {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	for k := range Accounts {
		token := oauth1.NewToken(Accounts[k].tok, Accounts[k].sec)
		Accounts[k].client = AuthConfig.Client(oauth1.NoContext, token)
		select {
		case <-ctx.Done():
			logger.Error(`Something wrong happened`,
				context.Cause(ctx))
			return
		default:
		}
		go Tweet(&Accounts[k], ctx, cancel)
	}

	select {
	case <-ctx.Done():
		if context.Cause(ctx) != nil {
			logger.Error(`Something wrong happened while reading from DMs`,
				context.Cause(ctx))
			return
		}
	}
}
