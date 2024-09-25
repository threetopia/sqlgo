package sqlgo

import "fmt"

type SQLGoGroupBy interface {
	SQLGroupBy(values ...sqlGoGroupByValue) SQLGoGroupBy
	SetSQLGroupBy(value sqlGoGroupByValue) SQLGoGroupBy

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoGroupBy
	SQLGoBase
}

type (
	sqlGoGroupByValue  string
	sqlGoGroupByValues []sqlGoGroupByValue

	sqlGoGroupBy struct {
		values         sqlGoGroupByValues
		sqlGoParameter SQLGoParameter
	}
)

func NewSQLGoGroupBy() SQLGoGroupBy {
	return new(sqlGoGroupBy)
}

func (s *sqlGoGroupBy) SQLGroupBy(values ...sqlGoGroupByValue) SQLGoGroupBy {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoGroupBy) SetSQLGroupBy(value sqlGoGroupByValue) SQLGoGroupBy {
	s.values = append(s.values, value)
	return s
}

func (s *sqlGoGroupBy) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoGroupBy {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoGroupBy) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoGroupBy) BuildSQL() string {
	var sql string
	if len(s.values) < 1 {
		return sql
	}
	sql = "GROUP BY "
	for i, v := range s.values {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}

		sql = fmt.Sprintf("%s%s", sql, v)
	}
	return sql
}
