package rpc

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/latifrons/lbserver/model"
)

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

func DebugBindJson(c *gin.Context, output interface{}) (err error) {
	err = c.ShouldBindJSON(output)
	return err
}

func ExtractPagingQuery(c *gin.Context) model.PagingParams {
	page := tryParseIntDefault(c.DefaultQuery("page", "1"), 1)
	size := tryParseIntDefault(c.DefaultQuery("size", "10"), 10)

	return model.PagingParams{
		Offset:    (page - 1) * size,
		Limit:     size,
		NeedTotal: true,
	}
}
