// Package main implements the hsm command
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/borud/nethsm"
)

var opt struct {
	APIURL   string `kong:"name='api',help='NetHSM REST Endpoint URL',default='https://127.0.0.1:8443/api/v1'"`
	NS       string `kong:"name='ns',help='Namespace'"`
	Username string `kong:"name='user',help='Username'"`
	Password string `kong:"name='pass',help='Password'"`
	SkipTLS  bool   `kong:"name='skip-tls',help='skip TLS verification',default=true"` // remove default
	// Debug     bool   `kong:"name='debug',help='turn on HTTPS debugging'"`

	Info      infoCmd      `kong:"cmd,name='info',help='get information on NetHSM instance'"`
	Health    healthCmd    `kong:"cmd,name='health',help='NetHSM health commands'"`
	Lock      lockCmd      `kong:"cmd,name='lock',help='lock NetHSM instance'"`
	Unlock    unlockCmd    `kong:"cmd,name='unlock',help='unlock NetHSM instance'"`
	Provision provisionCmd `kong:"cmd,name='provision',help='Provision NetHSM instance'"`
	User      userCmd      `kong:"cmd,name='user',help='user commands'"`
	Namespace namespaceCmd `kong:"cmd,name='ns',help='namespace commands'"`
}

func main() {
	ktx := kong.Parse(&opt,
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary:        false,
			Summary:             true,
			Compact:             true,
			Tree:                true,
			FlagsLast:           false,
			NoExpandSubcommands: false,
		}),
	)

	err := ktx.Run()
	if err != nil {
		fmt.Printf("\nError: %v\n\n", err)
	}
}

func withAuthClient(f func(ctx context.Context, client *nethsm.APIClient) error) error {
	// make it possible to skip TLS
	httpClient := http.DefaultClient
	if opt.SkipTLS {
		httpClient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	}

	// set the API endpoint URL and the http client
	config := nethsm.NewConfiguration()
	config.Servers = nethsm.ServerConfigurations{{URL: opt.APIURL}}
	config.HTTPClient = httpClient

	client := nethsm.NewAPIClient(config)

	ctx := context.WithValue(context.Background(), nethsm.ContextBasicAuth, nethsm.BasicAuth{
		UserName: opt.Username,
		Password: opt.Password,
	})

	return f(ctx, client)
}

func withClient(f func(ctx context.Context, client *nethsm.APIClient) error) error {
	// make it possible to skip TLS
	httpClient := http.DefaultClient
	if opt.SkipTLS {
		httpClient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	}

	// set the API endpoint URL and the http client
	config := nethsm.NewConfiguration()
	config.Servers = nethsm.ServerConfigurations{{URL: opt.APIURL}}
	config.HTTPClient = httpClient

	client := nethsm.NewAPIClient(config)
	ctx := context.Background()

	return f(ctx, client)
}
