package model

import (
	"github.com/latifrons/lbserver/berror"
	"gorm.io/gorm/clause"
	"strings"
)

const OrderDirectionASC = "ASC"
const OrderDirectionDESC = "DESC"

type PagingParams struct {
	Offset    int
	Limit     int
	NeedTotal bool
}

func (p PagingParams) ToPageNumSize() (pageNum int, pageSize int) {
	pageNum = int(p.Offset/p.Limit) + 1
	pageSize = p.Limit
	return
}

type PagingResult struct {
	Offset int
	Limit  int
	Total  int64
}

type OrderParams struct {
	OrderBy        string
	OrderDirection string
}

func (o *OrderParams) ToSqlOrderBy() (c interface{}, err error) {
	if o.OrderBy == "" || o.OrderDirection == "" {
		c = ""
		return
	}
	dir := strings.ToUpper(o.OrderDirection)
	if dir != OrderDirectionASC && dir != OrderDirectionDESC {
		err = berror.New(berror.ErrBadRequest, "bad order params", berror.CategoryBusinessFail)
		return
	}
	return clause.OrderByColumn{Column: clause.Column{Name: o.OrderBy}, Desc: dir == OrderDirectionDESC}, nil
}
