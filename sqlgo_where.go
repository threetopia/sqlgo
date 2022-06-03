package sqlgo

import "fmt"

var specialOperator = map[string]string{
	"ANY": "= ANY ",
	"IN":  " IN ",
}

type SQLGoWhere struct {
	values      []sqlGoWhereValues
	params      []interface{}
	paramCount  int
	isJoinScope bool
}

type sqlGoWhereValues struct {
	whereType   string
	whereColumn string
	operator    string
	value       interface{}
}

func NewSQLGoWhere() *SQLGoWhere {
	return &SQLGoWhere{}
}

func SetWhere(whereType string, whereColumn string, operator string, value interface{}) sqlGoWhereValues {
	return sqlGoWhereValues{
		whereType:   whereType,
		whereColumn: whereColumn,
		operator:    operator,
		value:       value,
	}
}

func (sw *SQLGoWhere) SQLWhere(values ...sqlGoWhereValues) *SQLGoWhere {
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
			sql = fmt.Sprintf("%s %s ", sql, v.whereType)
		}

		operator := v.operator
		if vo, ok := specialOperator[v.operator]; ok {
			v.operator = vo
		}

		switch vType := v.value.(type) {
		case *SQLGo:
			sql = fmt.Sprintf("%s%s%s(%s)", sql, v.whereColumn, v.operator, vType.SetParamsCount(sw.GetParamsCount()).BuildSQL())
			sw.SetParams(vType.GetParams()...)
			sw.SetParamsCount(vType.GetParamsCount())
		case []string:
			sql = buildWhereSlice(sw, sql, operator, v.whereType, v, vType)
		case []int:
			sql = buildWhereSlice(sw, sql, operator, v.whereType, v, vType)
		case []int64:
			sql = buildWhereSlice(sw, sql, operator, v.whereType, v, vType)
		case []float64:
			sql = buildWhereSlice(sw, sql, operator, v.whereType, v, vType)
		default:
			sw.SetParams(vType)
			sw.SetParamsCount(sw.GetParamsCount() + 1)
			sql = fmt.Sprintf("%s%s%s$%d", sql, v.whereColumn, v.operator, sw.GetParamsCount())
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

func (sw *SQLGoWhere) setJoinScope() *SQLGoWhere {
	sw.isJoinScope = true
	return sw
}

func buildWhereSlice[V string | int | int64 | float32 | float64](sw *SQLGoWhere, sql string, operator string, whereType string, vWhere sqlGoWhereValues, vType []V) string {
	if operator == "IN" {
		sql = fmt.Sprintf("%s%s%s%s(", sql, whereType, vWhere.whereColumn, vWhere.operator)
		for iIn, vIn := range vType {
			delimiter := ""
			if iIn > 0 {
				delimiter = ","
			}
			if sw.isJoinScope {
				sql = fmt.Sprintf("%s%s%x", sql, delimiter, vIn)
			} else {
				sw.SetParams(vIn)
				sw.SetParamsCount(sw.GetParamsCount() + 1)
				sql = fmt.Sprintf("%s%s$%d", sql, delimiter, sw.GetParamsCount())
			}
		}
		sql = fmt.Sprintf("%s)", sql)
	} else {
		sw.SetParams(vType)
		sw.SetParamsCount(sw.GetParamsCount() + 1)
		sql = fmt.Sprintf("%s%s%s%s($%d)", sql, whereType, vWhere.whereColumn, vWhere.operator, sw.GetParamsCount())
	}
	return sql
}
