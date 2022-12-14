package sqlgo

import "fmt"

type SQLGoFrom interface {
	SQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGoFrom

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoFrom
	SQLGoMandatory
}

type sqlGoFrom struct {
	table          sqlGoTable
	alias          sqlGoAlias
	sqlGoParameter SQLGoParameter
}

func NewSQLGoFrom() SQLGoFrom {
	return &sqlGoFrom{
		sqlGoParameter: NewSQLGoParameter(),
	}
}

func (s *sqlGoFrom) SQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGoFrom {
	s.table = table
	s.alias = alias
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

	if s.alias != "" {
		sql = fmt.Sprintf("%s AS %s", sql, s.alias)
	}
	return sql
}
