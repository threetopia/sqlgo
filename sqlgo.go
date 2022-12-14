package sqlgo

import "fmt"

type SQLGoMandatoryMethod interface {
	BuildSQL() string
}

type SQLGo interface {
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGo

	SQLSelect(values ...sqlGoSelectValue) SQLGo
	SetSQLSelect(alias sqlGoAlias, value interface{}) SQLGo

	SQLGoMandatoryMethod
}

type (
	sqlGo struct {
		sqlGoSelect    SQLGoSelect
		sqlGoParameter SQLGoParameter
	}

	sqlGoAlias string
	sqlGoValue interface{}
)

func NewSQLGo() SQLGo {
	return &sqlGo{
		sqlGoSelect:    NewSQLGoSelect(),
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

func (s *sqlGo) BuildSQL() string {
	sql := ""
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoSelect.BuildSQL())
	return sql
}
