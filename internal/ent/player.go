// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/msrevive/nexus2/internal/ent/player"
)

// Player is the model entity for the Player schema.
type Player struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Steamid holds the value of the "steamid" field.
	Steamid string `json:"steamid,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PlayerQuery when eager-loading is set.
	Edges PlayerEdges `json:"edges"`
}

// PlayerEdges holds the relations/edges for other nodes in the graph.
type PlayerEdges struct {
	// Characters holds the value of the characters edge.
	Characters []*Character `json:"characters,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// CharactersOrErr returns the Characters value or an error if the edge
// was not loaded in eager-loading.
func (e PlayerEdges) CharactersOrErr() ([]*Character, error) {
	if e.loadedTypes[0] {
		return e.Characters, nil
	}
	return nil, &NotLoadedError{edge: "characters"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Player) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case player.FieldSteamid:
			values[i] = new(sql.NullString)
		case player.FieldID:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Player", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Player fields.
func (pl *Player) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case player.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pl.ID = *value
			}
		case player.FieldSteamid:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field steamid", values[i])
			} else if value.Valid {
				pl.Steamid = value.String
			}
		}
	}
	return nil
}

// QueryCharacters queries the "characters" edge of the Player entity.
func (pl *Player) QueryCharacters() *CharacterQuery {
	return (&PlayerClient{config: pl.config}).QueryCharacters(pl)
}

// Update returns a builder for updating this Player.
// Note that you need to call Player.Unwrap() before calling this method if this Player
// was returned from a transaction, and the transaction was committed or rolled back.
func (pl *Player) Update() *PlayerUpdateOne {
	return (&PlayerClient{config: pl.config}).UpdateOne(pl)
}

// Unwrap unwraps the Player entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pl *Player) Unwrap() *Player {
	tx, ok := pl.config.driver.(*txDriver)
	if !ok {
		panic("ent: Player is not a transactional entity")
	}
	pl.config.driver = tx.drv
	return pl
}

// String implements the fmt.Stringer.
func (pl *Player) String() string {
	var builder strings.Builder
	builder.WriteString("Player(")
	builder.WriteString(fmt.Sprintf("id=%v", pl.ID))
	builder.WriteString(", steamid=")
	builder.WriteString(pl.Steamid)
	builder.WriteByte(')')
	return builder.String()
}

// Players is a parsable slice of Player.
type Players []*Player

func (pl Players) config(cfg config) {
	for _i := range pl {
		pl[_i].config = cfg
	}
}
