package main

import (
	"context"
	"fmt"

	"github.com/borud/nethsm"
)

type infoCmd struct{}

func (n *infoCmd) Run() error {
	return withClient(func(ctx context.Context, client *nethsm.APIClient) error {
		res, _, err := client.DefaultAPI.InfoGet(ctx).Execute()
		if err != nil {
			return err
		}

		fmt.Printf("Vendor: %s\nProduct: %s\n", res.Vendor, res.Product)
		return nil
	})
}
