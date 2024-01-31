package sqlgo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

type SQLGoWhere interface {
	SQLWhere(values ...sqlGoWhereValue) SQLGoWhere
	SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGoWhere
	SetSQLWhereBetween(whereType string, whereColumn string, firstVal, secondVal interface{}) SQLGoWhere
	SQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGoWhere
	SetSQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGoWhere

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoWhere
	SQLGoBase
}

type (
	sqlGoWhere struct {
		groupValue     sqlGoWhereGroupValueSlice
		sqlGOParameter SQLGoParameter
	}

	sqlGoWhereValue struct {
		whereType   string
		whereColumn string
		operator    string
		value       interface{}
		isParam     bool
	}

	sqlGoWhereGroupValue struct {
		valueSlice sqlGoWhereValueSlice
		whereType  string
	}

	sqlGoWhereBetween struct {
		firstVal  interface{}
		secondVal interface{}
	}

	sqlGoWhereGroupValueSlice []sqlGoWhereGroupValue
	sqlGoWhereValueSlice      []sqlGoWhereValue
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
	return new(sqlGoWhere)
}

func (s *sqlGoWhereValueSlice) Append(sql sqlGoWhereValue) {
	*s = append(*s, sql)
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

func SetSQLWhereBetween(whereType string, whereColumn string, firstVal, secondVal interface{}) sqlGoWhereValue {
	return sqlGoWhereValue{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    "BETWEEN",
		value:       sqlGoWhereBetween{firstVal: firstVal, secondVal: secondVal},
		isParam:     true,
	}
}

func SetSQLWheres(values ...sqlGoWhereValue) sqlGoWhereValueSlice {
	var wheres sqlGoWhereValueSlice
	for _, value := range values {
		wheres = append(wheres, value)
	}
	return wheres
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

func (s *sqlGoWhere) SQLWhere(valueSlice ...sqlGoWhereValue) SQLGoWhere {
	if len(s.groupValue) > 0 {
		s.groupValue[0].valueSlice = append(s.groupValue[0].valueSlice, valueSlice...)
	} else {
		s.groupValue = make(sqlGoWhereGroupValueSlice, 0)
		s.groupValue = append(s.groupValue, sqlGoWhereGroupValue{whereType: "AND", valueSlice: valueSlice})
	}
	return s
}

func (s *sqlGoWhere) SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGoWhere {
	if len(s.groupValue) > 0 {
		s.groupValue[0].valueSlice = append(s.groupValue[0].valueSlice, SetSQLWhere(whereType, whereColumn, operator, value))
	} else {
		s.groupValue = make(sqlGoWhereGroupValueSlice, 0)
		s.groupValue = append(s.groupValue, sqlGoWhereGroupValue{
			whereType:  "AND",
			valueSlice: append(make(sqlGoWhereValueSlice, 0), SetSQLWhere(whereType, whereColumn, operator, value)),
		})
	}
	return s
}

func (s *sqlGoWhere) SetSQLWhereBetween(whereType string, whereColumn string, firstVal, secondVal interface{}) SQLGoWhere {
	operator := "BETWEEN"
	value := sqlGoWhereBetween{firstVal: firstVal, secondVal: secondVal}
	if len(s.groupValue) > 0 {
		s.groupValue[0].valueSlice = append(s.groupValue[0].valueSlice, SetSQLWhere(whereType, whereColumn, operator, value))
	} else {
		s.groupValue = make(sqlGoWhereGroupValueSlice, 0)
		s.groupValue = append(s.groupValue, sqlGoWhereGroupValue{
			whereType:  "AND",
			valueSlice: append(make(sqlGoWhereValueSlice, 0), SetSQLWhere(whereType, whereColumn, operator, value)),
		})
	}
	return s
}

func (s *sqlGoWhere) SQLWhereGroup(whereType string, valueSlice ...sqlGoWhereValue) SQLGoWhere {
	if len(s.groupValue) < 1 {
		s.groupValue = make(sqlGoWhereGroupValueSlice, 0)
	}
	s.groupValue = append(s.groupValue, sqlGoWhereGroupValue{whereType: whereType, valueSlice: valueSlice})
	return s
}

func (s *sqlGoWhere) SetSQLWhereGroup(whereType string, valueSlice ...sqlGoWhereValue) SQLGoWhere {
	s.SQLWhereGroup(whereType, valueSlice...)
	return s
}

func (s *sqlGoWhere) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoWhere {
	s.sqlGOParameter = sqlGoParameter
	return s
}

func (s *sqlGoWhere) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGOParameter
}

func (s *sqlGoWhere) BuildSQL() string {
	var sql string
	if len(s.groupValue) < 1 {
		return sql
	}

	sql = "WHERE "
	for ig, vg := range s.groupValue {
		var sqlVal string
		for i, v := range vg.valueSlice {
			if i > 0 {
				sqlVal = fmt.Sprintf("%s %s ", sqlVal, strings.ToUpper(v.whereType))
			}

			operator := strings.ToUpper(v.operator)
			if vo, ok := specialOperator[operator]; ok {
				v.operator = vo
			}

			switch vType := v.value.(type) {
			case SQLGo:
				vType.SetSQLGoParameter(s.GetSQLGoParameter())
				sqlVal = fmt.Sprintf("%s%s%s(%s)", sqlVal, v.whereColumn, v.operator, vType.BuildSQL())
				s.SetSQLGoParameter(vType.GetSQLGoParameter())
			case sqlGoWhereBetween:
				sqlVal = buildWhereBetween(s, sqlVal, operator, v, vType)
			case []string:
				sqlVal = buildWhereSlice(s, sqlVal, operator, v, vType)
			case []int:
				sqlVal = buildWhereSlice(s, sqlVal, operator, v, vType)
			case []int64:
				sqlVal = buildWhereSlice(s, sqlVal, operator, v, vType)
			case []float64:
				sqlVal = buildWhereSlice(s, sqlVal, operator, v, vType)
			default:
				if !v.isParam {
					sqlVal = fmt.Sprintf("%s%s%s%s", sqlVal, v.whereColumn, v.operator, vType)
				} else {
					s.sqlGOParameter.SetSQLParameter(vType)
					sqlVal = fmt.Sprintf("%s%s%s%s", sqlVal, v.whereColumn, v.operator, s.GetSQLGoParameter().GetSQLParameterSign(vType))
				}
			}
		}
		if len(s.groupValue) > 1 {
			if ig > 0 {
				sql = fmt.Sprintf("%s %s ", sql, strings.ToUpper(vg.whereType))
			}
			sql = fmt.Sprintf("%s(%s)", sql, sqlVal)
		} else {
			sql = fmt.Sprintf("%s%s", sql, sqlVal)
		}
	}

	return sql
}

func buildWhereBetween(s SQLGoWhere, sql string, operator string, v sqlGoWhereValue, vType sqlGoWhereBetween) string {
	s.GetSQLGoParameter().SetSQLParameter(vType.firstVal)
	firstParamSign := s.GetSQLGoParameter().GetSQLParameterSign(vType.firstVal)
	s.GetSQLGoParameter().SetSQLParameter(vType.secondVal)
	secondParamSign := s.GetSQLGoParameter().GetSQLParameterSign(vType.secondVal)

	return fmt.Sprintf("%s(%s %s %s AND %s)", sql, v.whereColumn, v.operator, firstParamSign, secondParamSign)
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
