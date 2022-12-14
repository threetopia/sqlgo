package sqlgo

import "fmt"

const (
	Dialect           sqlGoDialect = "dialect"
	MySQLDialect      sqlGoDialect = "mysql-" + Dialect
	PostgreSQLDialect sqlGoDialect = "postgresql-" + Dialect
)

type SQLGoMandatory interface {
	GetSQLGoParameter() SQLGoParameter
	BuildSQL() string
}

type SQLGo interface {
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGo
	GetSQLGoParameter() SQLGoParameter

	SQLSelect(values ...sqlGoSelectValue) SQLGo
	SetSQLSelect(value interface{}, alias sqlGoAlias) SQLGo

	SQLInsert(table sqlGoTable, columns sqlGoInsertColumnList, values ...sqlGoInsertValueSlice) SQLGo
	SetSQLInsert(table sqlGoTable) SQLGo
	SetSQLInsertColumn(columns ...sqlGoInsertColumn) SQLGo
	SetSQLInsertValue(values ...sqlGoInsertValueSlice) SQLGo

	SQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGo

	SQLJoin(values ...sqlGoJoinValue) SQLGo
	SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) SQLGo

	SQLWhere(values ...sqlGoWhereValue) SQLGo
	SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGo

	SQLGoMandatory
}

type (
	sqlGo struct {
		sqlGoSelect    SQLGoSelect
		sqlGoInsert    SQLGoInsert
		sqlGoFrom      SQLGoFrom
		sqlGoJoin      SQLGoJoin
		sqlGoWhere     SQLGoWhere
		sqlGoParameter SQLGoParameter
	}

	sqlGoTable   interface{}
	sqlGoAlias   string
	sqlGoValue   interface{}
	sqlGoDialect string
)

func NewSQLGo() SQLGo {
	return &sqlGo{
		sqlGoSelect:    NewSQLGoSelect(),
		sqlGoInsert:    NewSQLGoInsert(),
		sqlGoFrom:      NewSQLGoFrom(),
		sqlGoJoin:      NewSQLGoJoin(),
		sqlGoWhere:     NewSQLGoWhere(),
		sqlGoParameter: NewSQLGoParameter(),
	}
}

func (s *sqlGo) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGo {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGo) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGo) SQLSelect(values ...sqlGoSelectValue) SQLGo {
	s.sqlGoSelect.SQLSelect(values...)
	return s
}

func (s *sqlGo) SetSQLSelect(value interface{}, alias sqlGoAlias) SQLGo {
	s.sqlGoSelect.SetSQLSelect(value, alias)
	return s
}

func (s *sqlGo) SQLInsert(table sqlGoTable, columns sqlGoInsertColumnList, values ...sqlGoInsertValueSlice) SQLGo {
	s.sqlGoInsert.SQLInsert(table, columns, values...)
	return s
}

func (s *sqlGo) SetSQLInsert(table sqlGoTable) SQLGo {
	s.sqlGoInsert.SetSQLInsert(table)
	return s
}

func (s *sqlGo) SetSQLInsertColumn(columns ...sqlGoInsertColumn) SQLGo {
	s.sqlGoInsert.SetSQLInsertColumn(columns...)
	return s
}

func (s *sqlGo) SetSQLInsertValue(values ...sqlGoInsertValueSlice) SQLGo {
	s.sqlGoInsert.SetSQLInsertValue(values...)
	return s
}

func (s *sqlGo) SQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGo {
	s.sqlGoFrom.SQLFrom(table, alias)
	return s
}

func (s *sqlGo) SQLJoin(values ...sqlGoJoinValue) SQLGo {
	s.sqlGoJoin.SQLJoin(values...)
	return s
}

func (s *sqlGo) SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) SQLGo {
	s.sqlGoJoin.SetSQLJoin(joinType, table, alias, sqlWhere...)
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
	s.sqlGoSelect.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoSelect.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoParameter.GetSQLGoParameter())

	s.sqlGoInsert.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoInsert.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoInsert.GetSQLGoParameter())

	s.sqlGoFrom.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoFrom.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoFrom.GetSQLGoParameter())

	s.sqlGoJoin.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoJoin.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoJoin.GetSQLGoParameter())

	s.sqlGoWhere.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoWhere.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoWhere.GetSQLGoParameter())
	return sql
}
