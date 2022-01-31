// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CharactersColumns holds the columns for the "characters" table.
	CharactersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "steamid", Type: field.TypeString},
		{Name: "slot", Type: field.TypeInt},
		{Name: "name", Type: field.TypeString},
		{Name: "gender", Type: field.TypeInt},
		{Name: "race", Type: field.TypeString},
		{Name: "flags", Type: field.TypeString, Default: "{}"},
		{Name: "quickslots", Type: field.TypeString, Default: "{}"},
		{Name: "quests", Type: field.TypeString, Default: "{}"},
		{Name: "guild", Type: field.TypeString},
		{Name: "kills", Type: field.TypeInt},
		{Name: "gold", Type: field.TypeInt},
		{Name: "skills", Type: field.TypeString, Default: "{}"},
		{Name: "pets", Type: field.TypeString, Default: "{}"},
		{Name: "health", Type: field.TypeInt},
		{Name: "mana", Type: field.TypeInt},
		{Name: "equipped", Type: field.TypeString, Default: "{}"},
		{Name: "lefthand", Type: field.TypeString},
		{Name: "righthand", Type: field.TypeString},
		{Name: "spells", Type: field.TypeString, Default: "{}"},
		{Name: "spellbook", Type: field.TypeString, Default: "{}"},
		{Name: "bags", Type: field.TypeString, Default: "{}"},
		{Name: "sheaths", Type: field.TypeString, Default: "{}"},
	}
	// CharactersTable holds the schema information for the "characters" table.
	CharactersTable = &schema.Table{
		Name:       "characters",
		Columns:    CharactersColumns,
		PrimaryKey: []*schema.Column{CharactersColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "character_id",
				Unique:  true,
				Columns: []*schema.Column{CharactersColumns[0]},
			},
			{
				Name:    "character_steamid_slot",
				Unique:  false,
				Columns: []*schema.Column{CharactersColumns[1], CharactersColumns[2]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CharactersTable,
	}
)

func init() {
}
