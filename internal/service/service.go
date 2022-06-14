package service

import (
	"context"

	"github.com/msrevive/nexus2/internal/ent"
)

func txn(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
