package postgres

import (
	"database/sql"
	"time"
)

func GetNullableInt32(value *int32) sql.NullInt32 {
	if value != nil {
		return sql.NullInt32{Int32: *value, Valid: true}
	}

	return sql.NullInt32{Valid: false}
}

func GetNullableString(value *string) sql.NullString {
	if value != nil {
		return sql.NullString{String: *value, Valid: true}
	}

	return sql.NullString{Valid: false}
}

func GetNullableTime(value *time.Time) sql.NullTime {
	if value != nil {
		return sql.NullTime{Time: *value, Valid: true}
	}

	return sql.NullTime{Valid: false}
}
