package sqlgo

import (
	"fmt"
)

type SQLGoOrder interface {
	SQLOrder(values ...sqlGoOrderValue) SQLGoOrder
	SetSQLOrder(value sqlGoValue, order string) SQLGoOrder
	SetSQLOrderEmbedding(column sqlGoColumn, operator sqlGoOperator, value sqlGoValue) SQLGoOrder

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoOrder
	SQLGoBase
}

type (
	sqlGoOrder struct {
		values         sqlGoOrderValues
		sqlGoParameter SQLGoParameter
	}

	sqlGoOrderValue struct {
		value sqlGoValue
		order string
	}

	sqlGoOrderValues []sqlGoOrderValue

	sqlGoOrderEmbedding struct {
		column   sqlGoColumn
		operator sqlGoOperator
		value    sqlGoValue
	}
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

func SetSQLOrderEmbedding(column sqlGoColumn, operator sqlGoOperator, value sqlGoValue) sqlGoOrderValue {
	return sqlGoOrderValue{
		value: sqlGoOrderEmbedding{
			column:   column,
			operator: operator,
			value:    value,
		},
	}
}

func (s *sqlGoOrder) SQLOrder(values ...sqlGoOrderValue) SQLGoOrder {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoOrder) SetSQLOrder(value sqlGoValue, order string) SQLGoOrder {
	s.values = append(s.values, SetSQLOrder(value, order))
	return s
}

func (s *sqlGoOrder) SetSQLOrderEmbedding(column sqlGoColumn, operator sqlGoOperator, value sqlGoValue) SQLGoOrder {
	s.values = append(s.values, SetSQLOrderEmbedding(column, operator, value))
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
	if len(s.values) < 1 {
		return sql
	}

	sql = "ORDER BY "
	for i, v := range s.values {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}

		switch vType := v.value.(type) {
		case SQLGo:
			vType.SetSQLGoParameter(s.GetSQLGoParameter())
			sql = fmt.Sprintf("(%s) %s", vType.BuildSQL(), v.order)
			s.SetSQLGoParameter(vType.GetSQLGoParameter())
		case sqlGoOrderEmbedding:
			s.sqlGoParameter.SetSQLParameter(vType)
			sql = fmt.Sprintf("%s%s %s %s", sql, vType.column, vType.operator, s.sqlGoParameter.GetSQLParameterSign(vType))
		default:
			sql = fmt.Sprintf("%s%s %s", sql, v.value, v.order)
		}
	}
	return sql
}
