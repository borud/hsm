package main

import (
	"context"
	"fmt"

	"github.com/borud/nethsm"
)

type healthCmd struct {
	State healthStateCmd `kong:"cmd,name='state',help='health state'"`
	Ready healthReadyCmd `kong:"cmd,name='ready',help='check if ready to take traffic (implies operational)'"`
	Alive healthAliveCmd `kong:"cmd,name='alive',help='check if alive but not ready to take traffic (implies locked or unprovisioned)'"`
}

type healthStateCmd struct{}
type healthReadyCmd struct{}
type healthAliveCmd struct{}

func (s *healthStateCmd) Run() error {
	return withClient(func(ctx context.Context, client *nethsm.APIClient) error {
		state, _, err := client.DefaultAPI.HealthStateGet(ctx).Execute()
		if err != nil {
			return fmt.Errorf("failed to get health state: %w", err)
		}
		fmt.Printf("%s\n", state.State)
		return nil
	})
}

func (r *healthReadyCmd) Run() error {
	return withClient(func(ctx context.Context, client *nethsm.APIClient) error {
		response, err := client.DefaultAPI.HealthReadyGet(ctx).Execute()
		if err != nil {
			if response.StatusCode == 412 {
				fmt.Println("no")
				return nil
			}
			return fmt.Errorf("failed to get ready state: %w", err)
		}

		if response.StatusCode == 200 {
			fmt.Println("yes")
			return nil
		}

		return fmt.Errorf("unknown response, code is %d", response.StatusCode)
	})
}
func (a *healthAliveCmd) Run() error {
	return withClient(func(ctx context.Context, client *nethsm.APIClient) error {
		response, err := client.DefaultAPI.HealthAliveGet(ctx).Execute()
		if err != nil {
			if response.StatusCode == 412 {
				fmt.Println("no")
				return nil
			}
			return fmt.Errorf("failed to get alive state: %w", err)
		}

		if response.StatusCode == 200 {
			fmt.Println("yes")
			return nil
		}

		return fmt.Errorf("unknown response, code is %d", response.StatusCode)
	})
}
