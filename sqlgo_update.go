package sqlgo

import "fmt"

type SQLGoUpdate struct {
	table      string
	values     []SQLGoUpdateValue
	params     []interface{}
	paramCount int
}

type SQLGoUpdateValue struct {
	column string
	value  interface{}
}

func NewSQLGoUpdate() *SQLGoUpdate {
	return &SQLGoUpdate{}
}

func SetUpdate(column string, value interface{}) SQLGoUpdateValue {
	return SQLGoUpdateValue{
		column: column,
		value:  value,
	}
}

func (su *SQLGoUpdate) SQLUpdate(table string, values ...SQLGoUpdateValue) *SQLGoUpdate {
	su.setSQLUpdateTable(table)
	su.setSQLUpdateValue(values...)
	return su
}

func (su *SQLGoUpdate) BuildSQL() string {
	if len(su.values) < 1 {
		return ""
	}

	sql := fmt.Sprintf("UPDATE %s SET (", su.table)
	for i, v := range su.values {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		su.SetParams(v.value)
		su.SetParamsCount(su.GetParamsCount() + 1)
		sql = fmt.Sprintf("%s%s=$%d", sql, v.column, su.GetParamsCount())
	}
	sql = fmt.Sprintf("%s)", sql)
	return sql
}

func (su *SQLGoUpdate) setSQLUpdateTable(table string) *SQLGoUpdate {
	su.table = table
	return su
}

func (su *SQLGoUpdate) setSQLUpdateValue(values ...SQLGoUpdateValue) *SQLGoUpdate {
	su.values = append(su.values, values...)
	return su
}

func (su *SQLGoUpdate) SetParams(params ...interface{}) *SQLGoUpdate {
	if len(params) < 1 {
		return su
	}
	su.params = append(su.params, params...)
	return su
}

func (su *SQLGoUpdate) GetParams() []interface{} {
	return su.params
}

func (su *SQLGoUpdate) SetParamsCount(paramsCount int) *SQLGoUpdate {
	su.paramCount = paramsCount
	return su
}

func (su *SQLGoUpdate) GetParamsCount() int {
	return su.paramCount
}
