package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/msrevive/nexus2/internal/ent"
	"github.com/msrevive/nexus2/internal/ent/character"
	"github.com/msrevive/nexus2/internal/ent/player"
)

type CharacterService interface {
	CharactersGetAll(ctx context.Context) ([]*ent.Character, error)
	CharactersGetBySteamid(ctx context.Context, sid string) ([]*ent.Character, error)
	CharacterGetBySteamidSlot(ctx context.Context, sid string, slt int) (*ent.Character, error)
	CharacterGetByID(ctx context.Context, id uuid.UUID) (*ent.Character, error)
	CharacterCreate(ctx context.Context, newChar ent.Character) (*ent.Character, error)
	CharacterUpdate(ctx context.Context, uid uuid.UUID, updateChar ent.Character) (*ent.Character, error)
	CharacterDelete(ctx context.Context, uid uuid.UUID) error
}

type charService struct {
	client *ent.Client
}

func NewCharacterService(client *ent.Client) *charService {
	return &charService{
		client: client,
	}
}

func (s *charService) CharactersGetAll(ctx context.Context) ([]*ent.Character, error) {
	chars, err := s.client.Character.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	return chars, nil
}

func (s *charService) CharactersGetBySteamid(ctx context.Context, sid string) ([]*ent.Character, error) {
	chars, err := s.client.Character.Query().Where(
		character.HasPlayerWith(player.Steamid(sid)),
	).All(ctx)
	if err != nil {
		return nil, err
	}

	return chars, nil
}

func (s *charService) CharacterGetBySteamidSlot(ctx context.Context, sid string, slt int) (*ent.Character, error) {
	char, err := s.client.Character.Query().Where(
		character.And(
			character.HasPlayerWith(player.Steamid(sid)),
			character.Slot(slt),
		),
	).Only(ctx)
	if err != nil {
		return nil, err
	}

	return char, nil
}

func (s *charService) CharacterGetByID(ctx context.Context, id uuid.UUID) (*ent.Character, error) {
	char, err := s.client.Character.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return char, nil
}

func (s *charService) CharacterCreate(ctx context.Context, sid string, newChar ent.Character) (*ent.Character, error) {
	var char *ent.Character
	err := txn(ctx, s.client, func(tx *ent.Tx) error {
		p, err := tx.Player.Query().Where(
			player.Steamid(sid),
		).Only(ctx)
		if err != nil {
			return err
		}

		c, err := tx.Character.Create().
			SetPlayer(p).
			SetSlot(newChar.Slot).
			SetSize(newChar.Size).
			SetData(newChar.Data).
			Save(ctx)
		if err != nil {
			return err
		}

		char = c
		return nil
	})

	if err != nil {
		return nil, err
	}

	return char, nil
}

func (s *charService) CharacterUpdate(ctx context.Context, uid uuid.UUID, updateChar ent.Character) (*ent.Character, error) {
	char, err := s.client.Character.UpdateOneID(uid).
		SetSize(updateChar.Size).
		SetData(updateChar.Data).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return char, nil
}

func (s *charService) CharacterDelete(ctx context.Context, uid uuid.UUID) error {
	err := s.client.Character.DeleteOneID(uid).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
