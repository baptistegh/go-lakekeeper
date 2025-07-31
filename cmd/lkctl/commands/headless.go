package commands

import (
	"context"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"

	log "github.com/sirupsen/logrus"

	"golang.org/x/oauth2/clientcredentials"
)

func MustCreateClient(ctx context.Context, opts clientOptions) *client.Client {

	opt := []client.ClientOptionFunc{}

	oauthConfig := clientcredentials.Config{
		ClientID:     opts.clientID,
		ClientSecret: opts.clientSecret,
		TokenURL:     opts.authURL,
		Scopes:       opts.scope,
	}

	if _, err := oauthConfig.Token(ctx); err != nil {
		log.Fatal(err)
	}

	as := core.OAuthTokenSource{
		TokenSource: oauthConfig.TokenSource(ctx),
	}

	if opts.boostrap {
		opt = append(opt, client.WithInitialBootstrapV1Enabled(true, true, core.Ptr(managementv1.ApplicationUserType)))
	}

	cli, err := client.NewAuthSourceClient(ctx, &as, opts.server, opt...)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}
