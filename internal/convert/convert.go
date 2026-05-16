package convert

import (
	"database/sql"

	"github.com/araddon/dateparse"
)

func DateStringToSqlNullTime(dateStr string) sql.NullTime {
    if dateStr == "" {
        return sql.NullTime{Valid: false}
    }

    t, err := dateparse.ParseAny(dateStr)
    if err != nil {
        return sql.NullTime{Valid: false}
    }

    return sql.NullTime{Time: t, Valid: true}
}