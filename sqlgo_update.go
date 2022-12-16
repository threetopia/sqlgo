package sqlgo

import "fmt"

type SQLGoUpdate interface {
	SQLUpdate(table sqlGoTable, values ...sqlGoUpdateValue) SQLGoUpdate
	SetSQLUpdate(table sqlGoTable) SQLGoUpdate
	SetSQLUpdateValue(values ...sqlGoUpdateValue) SQLGoUpdate

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoUpdate
	SQLGoMandatory
}

type (
	sqlGoUpdate struct {
		table          sqlGoTable
		values         []sqlGoUpdateValue
		sqlGoParameter SQLGoParameter
	}

	sqlGoUpdateValue struct {
		column string
		value  interface{}
	}
)

func NewSQLGoUpdate() SQLGoUpdate {
	return &sqlGoUpdate{
		sqlGoParameter: NewSQLGoParameter(),
	}
}

func SetSQLUpdate(column string, value interface{}) sqlGoUpdateValue {
	return sqlGoUpdateValue{
		column: column,
		value:  value,
	}
}

func (s *sqlGoUpdate) SQLUpdate(table sqlGoTable, values ...sqlGoUpdateValue) SQLGoUpdate {
	s.SetSQLUpdate(table)
	s.SetSQLUpdateValue(values...)
	return s
}

func (s *sqlGoUpdate) SetSQLUpdate(table sqlGoTable) SQLGoUpdate {
	s.table = table
	return s
}

func (s *sqlGoUpdate) SetSQLUpdateValue(values ...sqlGoUpdateValue) SQLGoUpdate {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoUpdate) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoUpdate {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoUpdate) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoUpdate) BuildSQL() string {
	var sql string
	if len(s.values) < 1 {
		return sql
	}

	sql = fmt.Sprintf("UPDATE %s SET ", s.table)
	for i, v := range s.values {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		s.sqlGoParameter.SetSQLParameter(v.value)
		sql = fmt.Sprintf("%s%s=%s", sql, v.column, s.sqlGoParameter.GetSQLParameterSign(v.value))
	}
	return sql
}
