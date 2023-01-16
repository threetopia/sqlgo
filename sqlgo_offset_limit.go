package sqlgo

import "fmt"

type SQLGoOffsetLimit interface {
	SQLOffsetLimit(offset int, limit int) SQLGoOffsetLimit
	SetSQLLimit(limit int) SQLGoOffsetLimit
	SetSQLOffset(offset int) SQLGoOffsetLimit
	SQLPageLimit(page int, limit int) SQLGoOffsetLimit
	SetSQLPage(page int) SQLGoOffsetLimit

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoOffsetLimit
	SQLGoMandatory
}

type sqlGoOffsetLimit struct {
	offset         int
	limit          int
	sqlGoParameter SQLGoParameter
}

func NewSQLGoOffsetLimit() SQLGoOffsetLimit {
	return new(sqlGoOffsetLimit)
}

func (s *sqlGoOffsetLimit) SQLOffsetLimit(offset int, limit int) SQLGoOffsetLimit {
	s.SetSQLOffset(offset)
	s.SetSQLLimit(limit)
	return s
}

func (s *sqlGoOffsetLimit) SetSQLLimit(limit int) SQLGoOffsetLimit {
	s.limit = limit
	return s
}

func (s *sqlGoOffsetLimit) SetSQLOffset(offset int) SQLGoOffsetLimit {
	s.offset = offset
	return s
}

func (s *sqlGoOffsetLimit) SQLPageLimit(page int, limit int) SQLGoOffsetLimit {
	s.SetSQLLimit(limit)
	s.SetSQLPage(page)
	return s
}

func (s *sqlGoOffsetLimit) SetSQLPage(page int) SQLGoOffsetLimit {
	if page > 0 {
		s.SetSQLOffset((page - 1) * s.limit)
	}
	return s
}

func (s *sqlGoOffsetLimit) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoOffsetLimit {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoOffsetLimit) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoOffsetLimit) BuildSQL() string {
	var sql string
	if s.offset >= 0 && s.limit > 0 {
		sql = fmt.Sprintf("OFFSET %d", s.offset)
	}
	if s.limit > 0 {
		if sql != "" {
			sql = fmt.Sprintf("%s ", sql)
		}
		sql = fmt.Sprintf("%sLIMIT %d", sql, s.limit)
	}
	return sql
}
