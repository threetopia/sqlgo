package sqlgo

import "fmt"

type SQLGoFrom interface {
	SQLFrom(table sqlGoTable, alias sqlGoAlias, columns ...sqlGoFromColumn) SQLGoFrom
	SetSQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGoFrom
	SetSQLFromColumn(columns ...sqlGoFromColumn) SQLGoFrom

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoFrom
	SQLGoMandatory
}

type (
	sqlGoFrom struct {
		table          sqlGoTable
		alias          sqlGoAlias
		columns        sqlGoFromColumnSlice
		sqlGoParameter SQLGoParameter
	}

	sqlGoFromColumn      string
	sqlGoFromColumnSlice []sqlGoFromColumn
)

func NewSQLGoFrom() SQLGoFrom {
	return &sqlGoFrom{
		sqlGoParameter: NewSQLGoParameter(),
	}
}

func (s *sqlGoFrom) SQLFrom(table sqlGoTable, alias sqlGoAlias, columns ...sqlGoFromColumn) SQLGoFrom {
	s.table = table
	s.alias = alias
	s.SetSQLFromColumn(columns...)
	return s
}

func (s *sqlGoFrom) SetSQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGoFrom {
	s.table = table
	s.alias = alias
	return s
}

func (s *sqlGoFrom) SetSQLFromColumn(columns ...sqlGoFromColumn) SQLGoFrom {
	s.columns = append(s.columns, columns...)
	return s
}

func (s *sqlGoFrom) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoFrom {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoFrom) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoFrom) BuildSQL() string {
	var sql string
	if s.table == nil {
		return sql
	}

	sql = " FROM "
	switch vType := s.table.(type) {
	case SQLGo:
		sql = fmt.Sprintf("%s(%s)", sql, vType.BuildSQL())
	default:
		sql = fmt.Sprintf("%s%s", sql, vType)
	}

	if len(s.columns) > 0 {
		sql = fmt.Sprintf("%s %s(", sql, s.alias)
		for i, v := range s.columns {
			if i > 0 {
				sql = fmt.Sprintf("%s, ", sql)
			}
			sql = fmt.Sprintf("%s%s", sql, v)
		}
		sql = fmt.Sprintf("%s)", sql)
	} else if s.alias != "" {
		sql = fmt.Sprintf("%s AS %s", sql, s.alias)
	}
	return sql
}
