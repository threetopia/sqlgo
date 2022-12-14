package sqlgo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

type SQLGoWhere interface {
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoWhere
	GetSQLGoParameter() SQLGoParameter

	SQLWhere(values ...sqlGoWhereValue) SQLGoWhere
	SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGoWhere

	SQLGoMandatory
}

type (
	sqlGoWhere struct {
		values         []sqlGoWhereValue
		sqlGOParameter SQLGoParameter
	}

	sqlGoWhereValue struct {
		whereType   string
		whereColumn string
		operator    string
		value       interface{}
		isParam     bool
	}
)

var specialOperator = map[string]string{
	"ANY":       "= ANY ",
	"ILIKE ANY": " ILIKE ANY ",
	"LIKE ANY":  " LIKE ANY ",
	"IN":        " IN ",
	"LIKE":      " LIKE ",
	"ILIKE":     " ILIKE ",
}

func NewSQLGoWhere() SQLGoWhere {
	return &sqlGoWhere{
		sqlGOParameter: NewSQLGoParameter(),
	}
}

func SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) sqlGoWhereValue {
	return sqlGoWhereValue{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    operator,
		value:       value,
		isParam:     true,
	}
}

func SetSQLWhereNotParam(whereType string, whereColumn string, operator string, value interface{}) sqlGoWhereValue {
	return sqlGoWhereValue{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    operator,
		value:       value,
		isParam:     false,
	}
}

func (s *sqlGoWhere) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoWhere {
	s.sqlGOParameter = sqlGoParameter
	return s
}

func (s *sqlGoWhere) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGOParameter
}

func (s *sqlGoWhere) SQLWhere(values ...sqlGoWhereValue) SQLGoWhere {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoWhere) SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGoWhere {
	s.values = append(s.values, SetSQLWhere(whereType, whereColumn, operator, value))
	return s
}

func (s *sqlGoWhere) BuildSQL() string {
	if len(s.values) < 1 {
		return ""
	}

	sql := "WHERE "
	for i, v := range s.values {
		if i > 0 {
			sql = fmt.Sprintf("%s %s ", sql, strings.ToUpper(v.whereType))
		}

		operator := strings.ToUpper(v.operator)
		if vo, ok := specialOperator[operator]; ok {
			v.operator = vo
		}

		switch vType := v.value.(type) {
		case SQLGo:
			vType.SetSQLGoParameter(s.GetSQLGoParameter())
			sql = fmt.Sprintf("%s%s%s(%s)", sql, v.whereColumn, v.operator, vType.BuildSQL())
			s.SetSQLGoParameter(vType.GetSQLGoParameter())
		case []string:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		case []int:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		case []int64:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		case []float64:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		default:
			if !v.isParam {
				sql = fmt.Sprintf("%s%s%s%s", sql, v.whereColumn, v.operator, vType)
			} else {
				s.sqlGOParameter.SetSQLParameter(vType)
				sql = fmt.Sprintf("%s%s%s%s", sql, v.whereColumn, v.operator, s.GetSQLGoParameter().GetSQLParameterSign(vType))
			}
		}
	}

	return sql
}

func buildWhereSlice[V string | int | int64 | float32 | float64](s SQLGoWhere, sql string, operator string, v sqlGoWhereValue, vType []V) string {
	loadedValue := make(map[V]bool)
	cleanVType := make([]V, 0)
	for _, vAny := range vType {
		if _, ok := loadedValue[vAny]; !ok {
			loadedValue[vAny] = true
		} else {
			continue
		}

		cleanVType = append(cleanVType, vAny)
	}
	vType = cleanVType

	if operator == "IN" {
		sql = fmt.Sprintf("%s%s%s(", sql, v.whereColumn, v.operator)
		for iIn, vIn := range vType {
			if iIn > 0 {
				sql = fmt.Sprintf("%s, ", sql)
			}

			if !v.isParam {
				sql = fmt.Sprintf("%s%x", sql, vIn)
			} else {
				s.GetSQLGoParameter().SetSQLParameter(vIn)
				sql = fmt.Sprintf("%s%s", sql, s.GetSQLGoParameter().GetSQLParameterSign(vIn))
			}
		}
		sql = fmt.Sprintf("%s)", sql)
	} else {
		if !v.isParam {
			sql = fmt.Sprintf("%s%s%s%x", sql, v.whereColumn, v.operator, vType)
		} else {
			var paramSign string
			if reflect.TypeOf(vType).Kind() == reflect.Slice {
				s.GetSQLGoParameter().SetSQLParameter(pq.Array(vType))
				paramSign = s.GetSQLGoParameter().GetSQLParameterSign(pq.Array(vType))
			} else {
				s.GetSQLGoParameter().SetSQLParameter(vType)
				paramSign = s.GetSQLGoParameter().GetSQLParameterSign(vType)
			}
			sql = fmt.Sprintf("%s%s%s(%s)", sql, v.whereColumn, v.operator, paramSign)
		}
	}
	return sql
}
