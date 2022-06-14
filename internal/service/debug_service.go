package service

import (
	"context"

	"github.com/msrevive/nexus2/internal/ent"
)

type DebugService interface {
	Debug(ctx context.Context) error
}

type dbgService struct {
	client *ent.Client
}

func NewDebugService(client *ent.Client) *dbgService {
	return &dbgService{
		client: client,
	}
}

func (s *dbgService) Debug(ctx context.Context) error {
	err := txn(ctx, s.client, func(tx *ent.Tx) error {
		player, err := tx.Player.Create().
			SetSteamid("76561198092541763").
			Save(ctx)
		if err != nil {
			return err
		}

		_, err = s.client.Character.Create().
			SetPlayer(player).
			SetSlot(1).
			SetData("data").
			Save(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
