// Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"context"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"

	log "github.com/sirupsen/logrus"

	"golang.org/x/oauth2/clientcredentials"
)

func MustCreateClient(ctx context.Context, opts *clientOptions) *client.Client {
	opt := []client.ClientOptionFunc{}

	switch {
	case opts.server == "":
		log.Fatal("You must provide server url")
	case opts.authURL == "":
		log.Fatal("You must provide auth url")
	case opts.clientID == "":
		log.Fatal("You must provide OAuth client_id")
	case opts.clientSecret == "":
		log.Fatal("You must provide OAuth client_secret")
	case len(opts.scope) == 0:
		log.Fatal("You must provide OAuth scope")
	}

	oauthConfig := clientcredentials.Config{
		ClientID:     opts.clientID,
		ClientSecret: opts.clientSecret,
		TokenURL:     opts.authURL,
		Scopes:       opts.scope,
	}

	log.Debug("testing OAuth2 client credentials")
	if _, err := oauthConfig.Token(ctx); err != nil {
		log.Fatal(err)
	}

	as := core.OAuthTokenSource{
		TokenSource: oauthConfig.TokenSource(ctx),
	}

	if opts.boostrap {
		log.Debug("enabling server bootstrap")
		opt = append(opt, client.WithInitialBootstrapV1Enabled(true, true, core.Ptr(managementv1.ApplicationUserType)))
	}

	cli, err := client.NewAuthSourceClient(ctx, &as, opts.server, opt...)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}
