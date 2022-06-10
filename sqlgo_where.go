package sqlgo

import (
	"fmt"
	"strings"
)

var specialOperator = map[string]string{
	"ANY":   "= ANY ",
	"IN":    " IN ",
	"LIKE":  " LIKE ",
	"ILIKE": " ILIKE ",
}

type SQLGoWhere struct {
	values     []SqlGoWhereValue
	params     []interface{}
	paramCount int
}

type SqlGoWhereValue struct {
	whereType   string
	whereColumn string
	operator    string
	value       interface{}
	isParam     bool
}

func NewSQLGoWhere() *SQLGoWhere {
	return &SQLGoWhere{}
}

func SetWhere(whereType string, whereColumn string, operator string, value interface{}) SqlGoWhereValue {
	return SqlGoWhereValue{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    operator,
		value:       value,
		isParam:     true,
	}
}

func SetWhereNotParam(whereType string, whereColumn string, operator string, value interface{}) SqlGoWhereValue {
	return SqlGoWhereValue{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    operator,
		value:       value,
		isParam:     false,
	}
}

func (sw *SQLGoWhere) SQLWhere(values ...SqlGoWhereValue) *SQLGoWhere {
	sw.values = append(sw.values, values...)
	return sw
}

func (sw *SQLGoWhere) BuildSQL() string {
	if len(sw.values) < 1 {
		return ""
	}

	sql := "WHERE "
	for i, v := range sw.values {
		if i > 0 {
			sql = fmt.Sprintf("%s %s ", sql, strings.ToUpper(v.whereType))
		}

		operator := strings.ToUpper(v.operator)
		if vo, ok := specialOperator[operator]; ok {
			v.operator = vo
		}

		switch vType := v.value.(type) {
		case *SQLGo:
			sql = fmt.Sprintf("%s%s%s(%s)", sql, v.whereColumn, v.operator, vType.SetParamsCount(sw.GetParamsCount()).BuildSQL())
			sw.SetParams(vType.GetParams()...).
				SetParamsCount(vType.GetParamsCount())
		case []string:
			sql = buildWhereSlice(sw, sql, operator, v, vType)
		case []int:
			sql = buildWhereSlice(sw, sql, operator, v, vType)
		case []int64:
			sql = buildWhereSlice(sw, sql, operator, v, vType)
		case []float64:
			sql = buildWhereSlice(sw, sql, operator, v, vType)
		default:
			if !v.isParam {
				sql = fmt.Sprintf("%s%s%s%s", sql, v.whereColumn, v.operator, vType)
			} else {
				sw.SetParams(vType)
				sw.SetParamsCount(sw.GetParamsCount() + 1)
				sql = fmt.Sprintf("%s%s%s$%d", sql, v.whereColumn, v.operator, sw.GetParamsCount())
			}
		}
	}

	return sql
}

func (sw *SQLGoWhere) SetParams(params ...interface{}) *SQLGoWhere {
	if len(params) < 1 {
		return sw
	}
	sw.params = append(sw.params, params...)
	return sw
}

func (sw *SQLGoWhere) GetParams() []interface{} {
	return sw.params
}

func (sw *SQLGoWhere) SetParamsCount(paramsCount int) *SQLGoWhere {
	sw.paramCount = paramsCount
	return sw
}

func (sw *SQLGoWhere) GetParamsCount() int {
	return sw.paramCount
}

func buildWhereSlice[V string | int | int64 | float32 | float64](sw *SQLGoWhere, sql string, operator string, v SqlGoWhereValue, vType []V) string {
	if operator == "IN" {
		sql = fmt.Sprintf("%s%s%s(", sql, v.whereColumn, v.operator)
		for iIn, vIn := range vType {
			if iIn > 0 {
				sql = fmt.Sprintf("%s, ", sql)
			}

			if !v.isParam {
				sql = fmt.Sprintf("%s%x", sql, vIn)
			} else {
				sw.SetParams(vIn)
				sw.SetParamsCount(sw.GetParamsCount() + 1)
				sql = fmt.Sprintf("%s$%d", sql, sw.GetParamsCount())
			}
		}
		sql = fmt.Sprintf("%s)", sql)
	} else {
		if !v.isParam {
			sql = fmt.Sprintf("%s%s%s%x", sql, v.whereColumn, v.operator, vType)
		} else {
			sw.SetParams(vType)
			sw.SetParamsCount(sw.GetParamsCount() + 1)
			sql = fmt.Sprintf("%s%s%s($%d)", sql, v.whereColumn, v.operator, sw.GetParamsCount())
		}
	}
	return sql
}
