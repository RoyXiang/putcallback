package putio

import (
	"context"
	"fmt"
	"log"

	put "github.com/putdotio/go-putio"
	"golang.org/x/oauth2"
)

func New(token string, maxTransfers int) *Put {
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(ctx, tokenSource)

	client := put.NewClient(oauthClient)
	info, err := client.Account.Info(ctx)
	if err != nil || !info.AccountActive {
		log.Fatal("You must have an active Put.io subscription")
	}
	if maxTransfers <= 0 || maxTransfers > info.SimultaneousDownloadLimit {
		maxTransfers = info.SimultaneousDownloadLimit
	}

	result := &Put{
		Client:                client,
		MaxTransfers:          maxTransfers,
		DefaultDownloadFolder: "",
	}
	if settings, err := client.Account.Settings(ctx); err == nil && settings.DefaultDownloadFolder != RootFolderId {
		fileInfo := result.GetFileInfo(settings.DefaultDownloadFolder)
		if fileInfo != nil {
			result.DefaultDownloadFolder = fmt.Sprintf("%s/", fileInfo.FullPath)
		}
	}
	return result
}
