package utils

import "github.com/jackc/pgx/v5/pgtype"

func PgTypeTextToString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

func StringToPgTypeText(t string) pgtype.Text {
	return pgtype.Text{
		String: t,
		Valid:  true,
	}
}
