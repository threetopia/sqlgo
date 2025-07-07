package sqlgo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

type SQLGoWhere interface {
	SQLWhere(values ...sqlGoWhereValue) SQLGoWhere
	SetSQLWhere(whereType, whereColumn, operator string, value sqlGoValue) SQLGoWhere
	SetSQLWhereBetween(whereType, whereColumn string, firstVal, secondVal sqlGoValue) SQLGoWhere
	SetSQLWhereToTsQuery(whereType, whereColumn, lang string, value sqlGoValue) SQLGoWhere
	SQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGoWhere
	SetSQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGoWhere

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoWhere
	SQLGoBase
}

type (
	sqlGoWhere struct {
		values sqlGoWhereValueSlice
		// groupValues    sqlGoWhereGroupValueSlice
		sqlGOParameter SQLGoParameter
	}

	sqlGoWhereValue struct {
		whereType   string
		whereColumn string
		operator    string
		value       sqlGoValue
		isParam     bool
	}
	sqlGoWhereValueSlice []sqlGoWhereValue

	sqlGoWhereBetween struct {
		firstVal  sqlGoValue
		secondVal sqlGoValue
	}

	sqlGoWhereToTsQuery struct {
		lang  string
		value sqlGoValue
	}
)

var specialOperator = map[string]string{
	"ANY":        "= ANY ",
	"ILIKE ANY":  " ILIKE ANY ",
	"LIKE ANY":   " LIKE ANY ",
	"IN":         " IN ",
	"NOT IN":     " NOT IN ",
	"LIKE":       " LIKE ",
	"NOT LIKE":   " NOT LIKE ",
	"ILIKE":      " ILIKE ",
	"NOT ILIKE":  " NOT ILIKE ",
	"TO TSQUERY": " TO TSQUERY ",
}

func NewSQLGoWhere() SQLGoWhere {
	return new(sqlGoWhere)
}

func (s *sqlGoWhereValueSlice) Append(sql sqlGoWhereValue) {
	*s = append(*s, sql)
}

func SetSQLWhere(whereType string, whereColumn string, operator string, value sqlGoValue) sqlGoWhereValue {
	return sqlGoWhereValue{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    operator,
		value:       value,
		isParam:     true,
	}
}

func SetSQLWhereGroup(whereType string, values ...sqlGoWhereValue) sqlGoWhereValue {
	return sqlGoWhereValue{
		whereType: whereType,
		value:     SetSQLWheres(values...),
	}
}

func SetSQLWhereBetween(whereType string, whereColumn string, firstVal, secondVal sqlGoValue) sqlGoWhereValue {
	return sqlGoWhereValue{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    "BETWEEN",
		value:       sqlGoWhereBetween{firstVal: firstVal, secondVal: secondVal},
		isParam:     true,
	}
}

func SetSQLWhereToTsQuery(whereType string, whereColumn string, lang string, value sqlGoValue) sqlGoWhereValue {
	return sqlGoWhereValue{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    "TS QUERY",
		value: sqlGoWhereToTsQuery{
			lang:  lang,
			value: value,
		},
		isParam: true,
	}
}

func SetSQLWheres(values ...sqlGoWhereValue) sqlGoWhereValueSlice {
	var wheres sqlGoWhereValueSlice
	for _, value := range values {
		wheres = append(wheres, value)
	}
	return wheres
}

func SetSQLWhereNotParam(whereType string, whereColumn string, operator string, value sqlGoValue) sqlGoWhereValue {
	return sqlGoWhereValue{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    operator,
		value:       value,
		isParam:     false,
	}
}

func (s *sqlGoWhere) SQLWhere(values ...sqlGoWhereValue) SQLGoWhere {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoWhere) SetSQLWhere(whereType string, whereColumn string, operator string, value sqlGoValue) SQLGoWhere {
	s.values = append(s.values, SetSQLWhere(whereType, whereColumn, operator, value))
	return s
}

func (s *sqlGoWhere) SetSQLWhereBetween(whereType string, whereColumn string, firstVal, secondVal sqlGoValue) SQLGoWhere {
	s.values = append(s.values, SetSQLWhereBetween(whereType, whereColumn, firstVal, secondVal))
	return s
}

func (s *sqlGoWhere) SetSQLWhereToTsQuery(whereType string, whereColumn string, lang string, value sqlGoValue) SQLGoWhere {
	s.values = append(s.values, SetSQLWhereToTsQuery(whereType, whereColumn, lang, value))
	return s
}

func (s *sqlGoWhere) SQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGoWhere {
	s.values = append(s.values, SetSQLWhereGroup(whereType, values...))
	return s
}

func (s *sqlGoWhere) SetSQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGoWhere {
	s.SQLWhereGroup(whereType, values...)
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
	if len(s.values) < 1 {
		return sql
	}

	sql = "WHERE "
	sql = fmt.Sprintf("%s%s", sql, buildWhereValues(s, s.values))

	return sql
}

func buildWhereValues(s SQLGoWhere, values sqlGoWhereValueSlice) string {
	var sql string
	for i, v := range values {
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
			s.SetSQLGoParameter(vType.GetSQLGoParameter())
		case sqlGoWhereToTsQuery:
			sql = buildWhereToTsQuery(s, sql, v, vType)
		case sqlGoWhereBetween:
			sql = buildWhereBetween(s, sql, v, vType)
		case sqlGoWhereValueSlice:
			sql = fmt.Sprintf("%s(%s)", sql, buildWhereValues(s, vType))
		case []string:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		case []int:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		case []int32:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		case []int64:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		case []float32:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		case []float64:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		case []bool:
			sql = buildWhereSlice(s, sql, operator, v, vType)
		default:
			if !v.isParam {
				sql = fmt.Sprintf("%s%s%s%s", sql, v.whereColumn, v.operator, vType)
			} else {
				s.GetSQLGoParameter().SetSQLParameter(vType)
				sql = fmt.Sprintf("%s%s%s%s", sql, v.whereColumn, v.operator, s.GetSQLGoParameter().GetSQLParameterSign(vType))
			}
		}
	}
	return sql
}

func buildWhereToTsQuery(s SQLGoWhere, sql string, v sqlGoWhereValue, value sqlGoWhereToTsQuery) string {
	var sqlVal string
	s.GetSQLGoParameter().SetSQLParameter(sqlVal)
	paramSign := s.GetSQLGoParameter().GetSQLParameterSign(sqlVal)
	return fmt.Sprintf("%s%s @@ to_tsquery('%s', %s)", sql, v.whereColumn, value.lang, paramSign)
}

func buildWhereBetween(s SQLGoWhere, sql string, v sqlGoWhereValue, vType sqlGoWhereBetween) string {
	s.GetSQLGoParameter().SetSQLParameter(vType.firstVal)
	firstParamSign := s.GetSQLGoParameter().GetSQLParameterSign(vType.firstVal)
	s.GetSQLGoParameter().SetSQLParameter(vType.secondVal)
	secondParamSign := s.GetSQLGoParameter().GetSQLParameterSign(vType.secondVal)

	return fmt.Sprintf("%s(%s %s %s AND %s)", sql, v.whereColumn, v.operator, firstParamSign, secondParamSign)
}

func buildWhereSlice[V string | int | int32 | int64 | float32 | float64 | bool](s SQLGoWhere, sql string, operator string, v sqlGoWhereValue, vType []V) string {
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

	switch operator {
	case "TS VECTOR":

	case "IN", "NOT IN":
		sql = fmt.Sprintf("%s%s%s(", sql, v.whereColumn, v.operator)
		for iIn, vIn := range vType {
			if iIn > 0 {
				sql = fmt.Sprintf("%s, ", sql)
			}

			if !v.isParam {
				sql = fmt.Sprintf("%s%v", sql, vIn)
			} else {
				s.GetSQLGoParameter().SetSQLParameter(vIn)
				sql = fmt.Sprintf("%s%s", sql, s.GetSQLGoParameter().GetSQLParameterSign(vIn))
			}
		}
		sql = fmt.Sprintf("%s)", sql)
	default:
		if !v.isParam {
			sql = fmt.Sprintf("%s%s%s%v", sql, v.whereColumn, v.operator, vType)
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
