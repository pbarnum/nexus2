// Code generated by entc, DO NOT EDIT.

package character

import (
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the character type in the database.
	Label = "character"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldSteamid holds the string denoting the steamid field in the database.
	FieldSteamid = "steamid"
	// FieldSlot holds the string denoting the slot field in the database.
	FieldSlot = "slot"
	// FieldSize holds the string denoting the size field in the database.
	FieldSize = "size"
	// FieldData holds the string denoting the data field in the database.
	FieldData = "data"
	// Table holds the table name of the character in the database.
	Table = "characters"
)

// Columns holds all SQL columns for character fields.
var Columns = []string{
	FieldID,
	FieldSteamid,
	FieldSlot,
	FieldSize,
	FieldData,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// SlotValidator is a validator for the "slot" field. It is called by the builders before save.
	SlotValidator func(int) error
	// DefaultSize holds the default value on creation for the "size" field.
	DefaultSize int
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
