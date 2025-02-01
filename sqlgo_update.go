package sqlgo

import (
	"fmt"
)

type SQLGoUpdate interface {
	SQLUpdate(table sqlGoTable, values ...sqlGoUpdateValue) SQLGoUpdate
	SetSQLUpdate(table sqlGoTable) SQLGoUpdate
	SetSQLUpdateValue(column string, value interface{}) SQLGoUpdate
	SetSQLUpdateValueSlice(values ...sqlGoUpdateValue) SQLGoUpdate

	SetSQLGoSchema(schema SQLGoSchema) SQLGoUpdate
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoUpdate
	SQLGoBase
}

type (
	sqlGoUpdate struct {
		table          sqlGoTable
		values         sqlGoUpdateValueSlice
		sqlGoSchema    SQLGoSchema
		sqlGoParameter SQLGoParameter
	}

	sqlGoUpdateValue struct {
		column string
		value  interface{}
	}
	sqlGoUpdateValueSlice []sqlGoUpdateValue

	sqlGoUpdateToTsVector struct {
		value interface{}
	}
)

func NewSQLGoUpdate() SQLGoUpdate {
	return new(sqlGoUpdate)
}

func SetSQLUpdateValue(column string, value interface{}) sqlGoUpdateValue {
	return sqlGoUpdateValue{
		column: column,
		value:  value,
	}
}

func SetSQLUpdateToTsVector(value interface{}) sqlGoUpdateToTsVector {
	return sqlGoUpdateToTsVector{
		value: value,
	}
}

func (s *sqlGoUpdate) SQLUpdate(table sqlGoTable, values ...sqlGoUpdateValue) SQLGoUpdate {
	s.SetSQLUpdate(table)
	s.SetSQLUpdateValueSlice(values...)
	return s
}

func (s *sqlGoUpdate) SetSQLUpdate(table sqlGoTable) SQLGoUpdate {
	s.table = table
	return s
}

func (s *sqlGoUpdate) SetSQLUpdateValueSlice(values ...sqlGoUpdateValue) SQLGoUpdate {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoUpdate) SetSQLUpdateValue(column string, value interface{}) SQLGoUpdate {
	s.SetSQLUpdateValueSlice(SetSQLUpdateValue(column, value))
	return s
}

func (s *sqlGoUpdate) SetSQLGoSchema(sqlGoSchema SQLGoSchema) SQLGoUpdate {
	s.sqlGoSchema = sqlGoSchema
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

	sql = fmt.Sprintf("UPDATE %s%s SET ", s.sqlGoSchema.BuildSQL(), s.table)
	for i, v := range s.values {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		switch vType := v.value.(type) {
		case sqlGoUpdateToTsVector:
			s.sqlGoParameter.SetSQLParameter(vType)
			sql = fmt.Sprintf("%s%s=%s", sql, v.column, s.sqlGoParameter.GetSQLParameterSign(vType))
		default:
			s.sqlGoParameter.SetSQLParameter(vType)
			sql = fmt.Sprintf("%s%s=%s", sql, v.column, s.sqlGoParameter.GetSQLParameterSign(vType))
		}
	}
	return sql
}
