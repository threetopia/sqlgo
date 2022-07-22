package sqlgo

import "fmt"

type SQLGoOffsetLimit struct {
	limit  int
	offset int
}

func NewSQLGoOffsetLimit() *SQLGoOffsetLimit {
	return &SQLGoOffsetLimit{
		offset: -1,
	}
}

func (slo *SQLGoOffsetLimit) SQLOffsetLimit(offset int, limit int) *SQLGoOffsetLimit {
	slo.limit = limit
	slo.offset = offset
	return slo
}

func (slo *SQLGoOffsetLimit) SetSQLLimit(limit int) *SQLGoOffsetLimit {
	slo.limit = limit
	return slo
}

func (slo *SQLGoOffsetLimit) SetSQLOffset(offset int) *SQLGoOffsetLimit {
	slo.offset = offset
	return slo
}

func (slo *SQLGoOffsetLimit) SQLPageLimit(page int, limit int) *SQLGoOffsetLimit {
	slo.limit = limit
	slo.offset = page * limit
	return slo
}

func (ss *SQLGoOffsetLimit) BuildSQL() string {
	sql := ""
	if ss.offset >= 0 {
		sql = fmt.Sprintf("%sOFFSET %d", sql, ss.offset)
	}
	if ss.limit > 0 {
		if sql != "" {
			sql = fmt.Sprintf("%s ", sql)
		}
		sql = fmt.Sprintf("%sLIMIT %d", sql, ss.limit)
	}
	return sql
}
