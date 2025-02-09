package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/borud/nethsm"
)

type lockCmd struct{}

type unlockCmd struct {
	UnlockPassphrase string `kong:"required,name='unlock-pass',help='unlock passphrase'"`
}

func (l *lockCmd) Run() error {
	return withAuthClient(func(ctx context.Context, client *nethsm.APIClient) error {
		_, err := client.DefaultAPI.LockPost(ctx).Execute()
		if err != nil {
			return fmt.Errorf("failed to lock NetHSM instance: %w", err)
		}
		slog.Info("successfully locked NetHSM instance", "apiURL", opt.APIURL)
		return nil
	})
}

func (u *unlockCmd) Run() error {
	return withClient(func(ctx context.Context, client *nethsm.APIClient) error {
		_, err := client.DefaultAPI.UnlockPost(ctx).
			UnlockRequestData(*nethsm.NewUnlockRequestData(u.UnlockPassphrase)).
			Execute()
		if err != nil {
			return fmt.Errorf("failed to unlock NetHSM instance: %w", err)
		}
		slog.Info("successfully unlocked NetHSM instance", "apiURL", opt.APIURL)
		return nil
	})
}
