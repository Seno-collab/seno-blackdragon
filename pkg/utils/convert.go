package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func StringFromPgText(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

func PgTextFromString(t string) pgtype.Text {
	return pgtype.Text{
		String: t,
		Valid:  true,
	}
}

// PgUUIDFromUUID converts google/uuid.UUID to pgtype.UUID
func PgUUIDFromUUID(u uuid.UUID) pgtype.UUID {
	var b [16]byte
	copy(b[:], u[:])
	return pgtype.UUID{
		Bytes: b,
		Valid: true,
	}
}

// UUIDFromPgUUID converts pgtype.UUID to google/uuid.UUID
// Returns zero uuid.UUID if not valid
func UUIDFromPgUUID(pg pgtype.UUID) uuid.UUID {
	if !pg.Valid {
		return uuid.UUID{}
	}
	return uuid.UUID(pg.Bytes)
}

// PgUUIDFromString parses a string UUID and converts it to pgtype.UUID
// Returns an invalid pgtype.UUID if parsing fails
func PgUUIDFromString(s string) pgtype.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return PgUUIDFromUUID(u)
}
