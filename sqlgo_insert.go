package sqlgo

import "fmt"

type SQLGOInsert struct {
	table      string
	columns    []SQLGoInsertColumn
	values     [][]SQLGoInsertValue
	params     []interface{}
	paramCount int
}

type SQLGoInsertColumn string

type SQLGoInsertValue interface{}

func NewSQLGOInsert() *SQLGOInsert {
	return &SQLGOInsert{}
}

func SetInsertColumns(columns ...SQLGoInsertColumn) []SQLGoInsertColumn {
	return columns
}

func SetInsertValues(values ...SQLGoInsertValue) []SQLGoInsertValue {
	return values
}

func (si *SQLGOInsert) SQLInsert(table string, columns []SQLGoInsertColumn, values ...[]SQLGoInsertValue) *SQLGOInsert {
	si.setSQLInsertTable(table)
	si.setSQLInsertColumn(columns...)
	si.setSQLInsertValue(values...)
	return si
}

func (si *SQLGOInsert) setSQLInsertTable(table string) *SQLGOInsert {
	si.table = table
	return si
}

func (si *SQLGOInsert) setSQLInsertColumn(columns ...SQLGoInsertColumn) *SQLGOInsert {
	si.columns = append(si.columns, columns...)
	return si
}

func (si *SQLGOInsert) setSQLInsertValue(values ...[]SQLGoInsertValue) *SQLGOInsert {
	si.values = append(si.values, values...)
	return si
}

func (si *SQLGOInsert) BuildSQL() string {
	if len(si.columns) < 1 {
		return ""
	}
	sql := fmt.Sprintf("INSERT %s (", si.table)
	for i, v := range si.columns {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		sql = fmt.Sprintf("%s%s", sql, v)
	}

	sql = fmt.Sprintf("%s)", sql)
	sql = fmt.Sprintf("%s VALUES ", sql)
	for iValues, vValues := range si.values {
		if iValues > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}

		sql = fmt.Sprintf("%s(", sql)
		for iValue, vValue := range vValues {
			if iValue > 0 {
				sql = fmt.Sprintf("%s, ", sql)
			}
			si.SetParams(vValue)
			si.SetParamsCount(si.GetParamsCount() + 1)
			sql = fmt.Sprintf("%s$%d", sql, si.GetParamsCount())
		}
		sql = fmt.Sprintf("%s)", sql)
	}
	return sql
}

func (si *SQLGOInsert) SetParams(params ...interface{}) *SQLGOInsert {
	if len(params) < 1 {
		return si
	}
	si.params = append(si.params, params...)
	return si
}

func (si *SQLGOInsert) GetParams() []interface{} {
	return si.params
}

func (si *SQLGOInsert) SetParamsCount(paramsCount int) *SQLGOInsert {
	si.paramCount = paramsCount
	return si
}

func (si *SQLGOInsert) GetParamsCount() int {
	return si.paramCount
}
