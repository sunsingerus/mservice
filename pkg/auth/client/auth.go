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

package client_auth

import (
	log "github.com/golang/glog"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"

	"context"
	cc "golang.org/x/oauth2/clientcredentials"

	"github.com/coreos/go-oidc"
)

func getOAuthClient(clientID, clientSecret, tokenURL string) (*oauth2.Token, error) {
	log.Infof("Setup OAuth params:\nClientID:%s\nClientSecret:%s\nTokenURL:%s\n", clientID, clientSecret, tokenURL)

	// client credential access
	config := &cc.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       []string{"profile"},
	}

	// you can modify the client (for example ignoring bad certs or otherwise)
	// by modifying the context
	ctx := context.Background()
	if token, err := config.Token(ctx); err == nil {
		log.Infof("Token received============\nAccessToken:\n%s\nTokenType:\n%s\nRefreshToken:\n%s\nExpiry:\n%s",
			token.AccessToken,
			token.TokenType,
			token.RefreshToken,
			token.Expiry,
		)
		return token, nil
	} else {
		log.Infof("Error token request %v", err)
		return nil, err
	}
}

func fetchToken(clientID, clientSecret, tokenURL string) *oauth2.Token {
	token, _ := getOAuthClient(clientID, clientSecret, tokenURL)
	return token

	//	return &oauth2.Token{
	//		AccessToken: "my-secret-token",
	//	}
}

func SetupOAuth(clientID, clientSecret, tokenURL string) ([]grpc.DialOption, error) {
	perRPCCredentials := oauth.NewOauthAccess(fetchToken(clientID, clientSecret, tokenURL))

	// Set token once per connection
	// It will be sent by gRPC on each call, without need to do it manually
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPCCredentials),
	}

	return opts, nil
}

func qwe() {
	_, err := oidc.NewProvider(context.Background(), "providerURI")
	if err != nil {
		log.Fatal(err)
	}
}
