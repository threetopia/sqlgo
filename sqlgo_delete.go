package sqlgo

import "fmt"

type SQLGoDelete interface {
	SQLDelete(table sqlGoTable) SQLGoDelete
	SetSQLDelete(table sqlGoTable) SQLGoDelete

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoDelete
	SQLGoMandatory
}

type sqlGoDelete struct {
	table          sqlGoTable
	sqlGoParameter SQLGoParameter
}

func NewSQLGoDelete() SQLGoDelete {
	return &sqlGoDelete{}
}

func (s *sqlGoDelete) SQLDelete(table sqlGoTable) SQLGoDelete {
	s.table = table
	return s
}

func (s *sqlGoDelete) SetSQLDelete(table sqlGoTable) SQLGoDelete {
	s.SQLDelete(table)
	return s
}

func (s *sqlGoDelete) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoDelete {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoDelete) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoDelete) BuildSQL() string {
	var sql string
	if s.table == nil {
		return sql
	}

	sql = fmt.Sprintf("DELETE FROM %s", s.table)
	return sql
}
