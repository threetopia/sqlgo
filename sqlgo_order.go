package sqlgo

import (
	"fmt"
)

type SQLGoOrder interface {
	SQLOrder(values ...sqlGoOrderValue) SQLGoOrder
	SetSQLOrder(value sqlGoValue, order string) SQLGoOrder

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoOrder
	SQLGoMandatory
}

type (
	sqlGoOrder struct {
		valueSlice     sqlGoOrderValueSlice
		sqlGoParameter SQLGoParameter
	}

	sqlGoOrderValue struct {
		value sqlGoValue
		order string
	}

	sqlGoOrderValueSlice []sqlGoOrderValue
)

func NewSQLGoOrder() SQLGoOrder {
	return new(sqlGoOrder)
}

func SetSQLOrder(value sqlGoValue, order string) sqlGoOrderValue {
	return sqlGoOrderValue{
		value: value,
		order: order,
	}
}

func (s *sqlGoOrder) SQLOrder(values ...sqlGoOrderValue) SQLGoOrder {
	s.valueSlice = append(s.valueSlice, values...)
	return s
}

func (s *sqlGoOrder) SetSQLOrder(value sqlGoValue, order string) SQLGoOrder {
	s.valueSlice = append(s.valueSlice, SetSQLOrder(value, order))
	return s
}

func (s *sqlGoOrder) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoOrder {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoOrder) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoOrder) BuildSQL() string {
	var sql string
	if len(s.valueSlice) < 1 {
		return sql
	}

	sql = "ORDER BY "
	for i, v := range s.valueSlice {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}

		switch vType := v.value.(type) {
		case SQLGo:
			vType.SetSQLGoParameter(s.GetSQLGoParameter())
			sql = fmt.Sprintf("(%s) %s", vType.BuildSQL(), v.order)
			s.SetSQLGoParameter(vType.GetSQLGoParameter())
		default:
			sql = fmt.Sprintf("%s%s %s", sql, v.value, v.order)
		}
	}

	return sql
}
