package sqlgo

import "fmt"

type SQLGoWhere struct {
	values     []sqlGoWhereValues
	params     []interface{}
	paramCount int
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
		switch vType := v.value.(type) {
		case string:
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
