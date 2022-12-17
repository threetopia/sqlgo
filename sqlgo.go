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
	SQLSelect(values ...sqlGoSelectValue) SQLGo
	SetSQLSelect(value interface{}, alias sqlGoAlias) SQLGo

	SQLInsert(table sqlGoTable, columns sqlGoInsertColumnSlice, values ...sqlGoInsertValueSlice) SQLGo
	SetSQLInsert(table sqlGoTable) SQLGo
	SetSQLInsertColumn(columns ...sqlGoInsertColumn) SQLGo
	SetSQLInsertValue(values ...sqlGoInsertValueSlice) SQLGo

	SQLUpdate(table sqlGoTable, values ...sqlGoUpdateValue) SQLGo
	SetSQLUpdate(table sqlGoTable) SQLGo
	SetSQLUpdateValue(values ...sqlGoUpdateValue) SQLGo

	SQLDelete(table sqlGoTable) SQLGo
	SetSQLDelete(table sqlGoTable) SQLGo

	SQLValues(values ...sqlGoValuesValueSlice) SQLGo
	SetSQLValues(values ...sqlGoValuesValueSlice) SQLGo

	SQLFrom(table sqlGoTable, alias sqlGoAlias, columns ...sqlGoFromColumn) SQLGo
	SetSQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGo
	SetSQLFromColumn(columns ...sqlGoFromColumn) SQLGo

	SQLJoin(values ...sqlGoJoinValue) SQLGo
	SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) SQLGo

	SQLWhere(values ...sqlGoWhereValue) SQLGo
	SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGo

	SQLOffsetLimit(offset int, limit int) SQLGo
	SetSQLLimit(limit int) SQLGo
	SetSQLOffset(offset int) SQLGo
	SQLPageLimit(page int, limit int) SQLGo
	SetSQLPage(page int) SQLGo

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGo
	SQLGoMandatory
}

type (
	sqlGo struct {
		sqlGoSelect      SQLGoSelect
		sqlGoInsert      SQLGoInsert
		sqlGoUpdate      SQLGoUpdate
		sqlGoDelete      SQLGoDelete
		sqlGoValues      SQLGoValues
		sqlGoFrom        SQLGoFrom
		sqlGoJoin        SQLGoJoin
		sqlGoWhere       SQLGoWhere
		sqlGoOffsetLimit SQLGoOffsetLimit
		sqlGoParameter   SQLGoParameter
	}

	sqlGoTable   interface{}
	sqlGoAlias   string
	sqlGoValue   interface{}
	sqlGoDialect string
)

func NewSQLGo() SQLGo {
	return &sqlGo{
		sqlGoSelect:      NewSQLGoSelect(),
		sqlGoInsert:      NewSQLGoInsert(),
		sqlGoUpdate:      NewSQLGoUpdate(),
		sqlGoDelete:      NewSQLGoDelete(),
		sqlGoValues:      NewSQLGoValues(),
		sqlGoFrom:        NewSQLGoFrom(),
		sqlGoJoin:        NewSQLGoJoin(),
		sqlGoWhere:       NewSQLGoWhere(),
		sqlGoOffsetLimit: NewSQLGoOffsetLimit(),
		sqlGoParameter:   NewSQLGoParameter(),
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

func (s *sqlGo) SQLInsert(table sqlGoTable, columns sqlGoInsertColumnSlice, values ...sqlGoInsertValueSlice) SQLGo {
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

func (s *sqlGo) SQLUpdate(table sqlGoTable, values ...sqlGoUpdateValue) SQLGo {
	s.sqlGoUpdate.SQLUpdate(table, values...)
	return s
}

func (s *sqlGo) SetSQLUpdate(table sqlGoTable) SQLGo {
	s.sqlGoUpdate.SetSQLUpdate(table)
	return s
}

func (s *sqlGo) SetSQLUpdateValue(values ...sqlGoUpdateValue) SQLGo {
	s.sqlGoUpdate.SetSQLUpdateValue(values...)
	return s
}

func (s *sqlGo) SQLDelete(table sqlGoTable) SQLGo {
	s.sqlGoDelete.SQLDelete(table)
	return s
}

func (s *sqlGo) SetSQLDelete(table sqlGoTable) SQLGo {
	s.sqlGoDelete.SQLDelete(table)
	return s
}

func (s *sqlGo) SQLValues(values ...sqlGoValuesValueSlice) SQLGo {
	s.sqlGoValues.SQLValues(values...)
	return s
}

func (s *sqlGo) SetSQLValues(values ...sqlGoValuesValueSlice) SQLGo {
	s.sqlGoValues.SetSQLValues(values...)
	return s
}

func (s *sqlGo) SQLFrom(table sqlGoTable, alias sqlGoAlias, columns ...sqlGoFromColumn) SQLGo {
	s.sqlGoFrom.SQLFrom(table, alias, columns...)
	return s
}

func (s *sqlGo) SetSQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGo {
	s.sqlGoFrom.SetSQLFrom(table, alias)
	return s
}

func (s *sqlGo) SetSQLFromColumn(columns ...sqlGoFromColumn) SQLGo {
	s.sqlGoFrom.SetSQLFromColumn(columns...)
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

func (s *sqlGo) SQLOffsetLimit(offset int, limit int) SQLGo {
	s.sqlGoOffsetLimit.SQLOffsetLimit(offset, limit)
	return s
}

func (s *sqlGo) SetSQLLimit(limit int) SQLGo {
	s.sqlGoOffsetLimit.SetSQLLimit(limit)
	return s
}

func (s *sqlGo) SetSQLOffset(offset int) SQLGo {
	s.sqlGoOffsetLimit.SetSQLOffset(offset)
	return s
}

func (s *sqlGo) SQLPageLimit(page int, limit int) SQLGo {
	s.sqlGoOffsetLimit.SQLPageLimit(page, limit)
	return s
}

func (s *sqlGo) SetSQLPage(page int) SQLGo {
	s.sqlGoOffsetLimit.SetSQLPage(page)
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

	s.sqlGoUpdate.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoUpdate.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoUpdate.GetSQLGoParameter())

	s.sqlGoDelete.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoDelete.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoDelete.GetSQLGoParameter())

	s.sqlGoFrom.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoFrom.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoFrom.GetSQLGoParameter())

	s.sqlGoJoin.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoJoin.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoJoin.GetSQLGoParameter())

	s.sqlGoWhere.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoWhere.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoWhere.GetSQLGoParameter())

	s.sqlGoValues.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoValues.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoValues.GetSQLGoParameter())

	s.sqlGoOffsetLimit.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = fmt.Sprintf("%s%s", sql, s.sqlGoOffsetLimit.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoOffsetLimit.GetSQLGoParameter())

	return sql
}
