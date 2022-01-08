package putio

import (
	"context"
	"log"

	put "github.com/putdotio/go-putio"
	"golang.org/x/oauth2"
)

func New(token string) *Put {
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(ctx, tokenSource)

	client := put.NewClient(oauthClient)
	info, err := client.Account.Info(ctx)
	if err != nil || !info.AccountActive {
		log.Fatal("You must have an active Put.io subscription")
	}

	return &Put{
		Client:       client,
		MaxTransfers: info.SimultaneousDownloadLimit,
	}
}
