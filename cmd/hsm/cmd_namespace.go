package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/borud/nethsm"
)

type namespaceCmd struct {
	Add    addNamespaceCmd    `kong:"cmd,name='add',help='add namespace'"`
	List   listNamespaceCmd   `kong:"cmd,name='ls',help='list namespaces'"`
	Delete deleteNamespaceCmd `kong:"cmd,name='del',help='delete namespace'"`
}

type addNamespaceCmd struct {
	ID string `kong:"required,name='id',help='id of namespace to be created'"`
}

type listNamespaceCmd struct{}

type deleteNamespaceCmd struct {
	ID string `kong:"required,name='id',help='id of namespace to be deleted'"`
}

func (a *addNamespaceCmd) Run() error {
	return withAuthClient(func(ctx context.Context, client *nethsm.APIClient) error {
		_, err := client.DefaultAPI.NamespacesNamespaceIDPut(ctx, a.ID).Execute()
		if err != nil {
			return fmt.Errorf("failed to create namespace [%s]: %w", a.ID, err)
		}

		slog.Info("created namespace", "id", a.ID)
		return nil
	})
}

func (l *listNamespaceCmd) Run() error {
	return withAuthClient(func(ctx context.Context, client *nethsm.APIClient) error {
		namespaces, _, err := client.DefaultAPI.NamespacesGet(ctx).Execute()
		if err != nil {
			return fmt.Errorf("failed to list namespaces: %w", err)
		}

		for _, nsitem := range namespaces {
			fmt.Println(nsitem.Id)
		}
		return nil
	})
}

func (d *deleteNamespaceCmd) Run() error {
	return withAuthClient(func(ctx context.Context, client *nethsm.APIClient) error {
		_, err := client.DefaultAPI.NamespacesNamespaceIDDelete(ctx, d.ID).Execute()
		if err != nil {
			return fmt.Errorf("failed to delete namespace [%s]: %w", d.ID, err)
		}
		slog.Info("deleted namespace", "id", d.ID)
		return nil
	})
}
