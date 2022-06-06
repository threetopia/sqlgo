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

func (ss *SQLGoOffsetLimit) BuildSQL() string {
	sql := ""
	if ss.offset >= 0 {
		sql = fmt.Sprintf("%s OFFSET %d", sql, ss.offset)
	}
	if ss.limit > 0 {
		sql = fmt.Sprintf("%s LIMIT %d", sql, ss.limit)
	}
	return sql
}
