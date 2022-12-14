package sqlgo

import "fmt"

type SQLGoMandatory interface {
	BuildSQL() string
}

type SQLGo interface {
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGo

	SQLSelect(values ...sqlGoSelectValue) SQLGo
	SetSQLSelect(alias sqlGoAlias, value interface{}) SQLGo

	SQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGo

	SQLWhere(values ...sqlGoWhereValue) SQLGo
	SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGo

	SQLGoMandatory
}

type (
	sqlGo struct {
		sqlGoSelect    SQLGoSelect
		sqlGoFrom      SQLGoFrom
		sqlGoWhere     SQLGoWhere
		sqlGoParameter SQLGoParameter
	}

	sqlGoTable interface{}
	sqlGoAlias string
	sqlGoValue interface{}
)

func NewSQLGo() SQLGo {
	return &sqlGo{
		sqlGoSelect:    NewSQLGoSelect(),
		sqlGoFrom:      NewSQLGoFrom(),
		sqlGoWhere:     NewSQLGoWhere(),
		sqlGoParameter: NewSQLGoParameter(),
	}
}

func (s *sqlGo) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGo {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGo) SQLSelect(values ...sqlGoSelectValue) SQLGo {
	s.sqlGoSelect.SQLSelect(values...)
	return s
}

func (s *sqlGo) SetSQLSelect(alias sqlGoAlias, value interface{}) SQLGo {
	s.sqlGoSelect.SetSQLSelect(alias, value)
	return s
}

func (s *sqlGo) SQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGo {
	s.sqlGoFrom.SQLFrom(table, alias)
	return s
}

func (s *sqlGo) SQLWhere(values ...sqlGoWhereValue) SQLGo {
	s.sqlGoWhere.SQLWhere(values...)
	return s
}

func (s *sqlGo) SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGo {
	s.sqlGoWhere.SetSQLWhere(whereType, whereColumn, operator, value)
	return s
}

func (s *sqlGo) BuildSQL() string {
	sql := ""
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoSelect.BuildSQL())
	sql = fmt.Sprintf("%s %s", sql, s.sqlGoFrom.BuildSQL())
	sql = fmt.Sprintf("%s %s", sql, s.sqlGoWhere.BuildSQL())
	return sql
}
