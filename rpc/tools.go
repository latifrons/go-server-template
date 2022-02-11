package rpc

import "database/sql"

func SqlNullTimeToInt64Default(value sql.NullTime) int64 {
	if value.Valid {
		return value.Time.Unix()
	}
	return 0
}

func SqlNullStringToStringDefault(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}
