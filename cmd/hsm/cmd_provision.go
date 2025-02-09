package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/borud/nethsm"
)

type provisionCmd struct {
	UnlockPass string `kong:"required,name='unlock-pass',help:'unlock passphrase'"`
	AdminPass  string `kong:"required,name='admin-pass',help:'admin passphrase'"`
}

func (p *provisionCmd) Run() error {
	return withClient(func(ctx context.Context, client *nethsm.APIClient) error {
		slog.Info("provisioning NetHSM, this may take some time", "apiURL", opt.APIURL)
		_, err := client.DefaultAPI.
			ProvisionPost(ctx).
			ProvisionRequestData(
				*nethsm.NewProvisionRequestData(
					p.UnlockPass,
					p.AdminPass,
					time.Time{})).Execute()
		if err != nil {
			return fmt.Errorf("failed to provision NetHSM instance: %w", err)
		}
		slog.Info("successfully provisioned NetHSM instance", "apiURL", opt.APIURL)
		return err
	})
}
