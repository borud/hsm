package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/borud/nethsm"
)

type userCmd struct {
	Add    addUserCmd    `kong:"cmd,name='add',help='add user'"`
	Get    getUserCmd    `kong:"cmd,name='get',help='get user'"`
	List   listUserCmd   `kong:"cmd,name='ls',help='list users'"`
	Delete deleteUserCmd `kong:"cmd,name='del',help='delete user'"`
}

type addUserCmd struct {
	ID         string `kong:"required,name='id',help='id of user to be created'"`
	RealName   string `kong:"required,name='realname',help='full name of user to be created'"`
	Role       string `kong:"required,name='role',help='user role'"`
	PassPhrase string `kong:"required,name='passphrase',help='passphrase for user'"`
}

type getUserCmd struct {
	ID string `kong:"required,name='id',help='id of user to be fetched'"`
}

type listUserCmd struct{}

type deleteUserCmd struct {
	ID string `kong:"required,name='id',help='id of user to be deleted'"`
}

func (a *addUserCmd) Run() error {
	return withAuthClient(func(ctx context.Context, client *nethsm.APIClient) error {
		// make sure hte role is valid
		role := nethsm.UserRole(a.Role)
		if !role.IsValid() {
			return fmt.Errorf("invalid user role [%s], valid roles are: %v", a.Role, nethsm.AllowedUserRoleEnumValues)
		}

		userPostData := nethsm.NewUserPostData(a.RealName, role, a.PassPhrase)
		res, err := client.DefaultAPI.UsersUserIDPut(ctx, a.ID).UserPostData(*userPostData).Execute()
		if err != nil {
			return fmt.Errorf("failed to create user [%s] res[%+v]: %w", a.ID, res, err)
		}

		slog.Info("created user", "id", a.ID)
		return nil
	})
}

func (g *getUserCmd) Run() error {
	return withAuthClient(func(ctx context.Context, client *nethsm.APIClient) error {
		userData, _, err := client.DefaultAPI.UsersUserIDGet(ctx, g.ID).Execute()
		if err != nil {
			return fmt.Errorf("failed to get user [%s]: %w", g.ID, err)
		}
		fmt.Printf("ID      : %s\nRealname: %s\nRole    : %s\n", g.ID, userData.RealName, userData.Role)
		return nil
	})
}

func (l *listUserCmd) Run() error {
	return withAuthClient(func(ctx context.Context, client *nethsm.APIClient) error {
		users, _, err := client.DefaultAPI.UsersGet(ctx).Execute()
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		for _, item := range users {
			fmt.Println(item.User)
		}
		return nil
	})
}

func (d *deleteUserCmd) Run() error {
	return withAuthClient(func(ctx context.Context, client *nethsm.APIClient) error {
		_, err := client.DefaultAPI.UsersUserIDDelete(ctx, d.ID).Execute()
		if err != nil {
			return fmt.Errorf("failed to delete user [%s]: %w", d.ID, err)
		}
		fmt.Printf("deleted user with id [%s]\n", d.ID)
		return nil
	})
}
