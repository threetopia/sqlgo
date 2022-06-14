package sqlgo

import "fmt"

type SQLGoDelete struct {
	table      string
	params     []interface{}
	paramCount int
}

func NewSQLGoDelete() *SQLGoDelete {
	return &SQLGoDelete{}
}

func (sd *SQLGoDelete) SQLDelete(table string) *SQLGoDelete {
	sd.table = table
	return sd
}

func (sd *SQLGoDelete) BuildSQL() string {
	if sd.table == "" {
		return ""
	}

	sql := fmt.Sprintf("DELETE FROM %s", sd.table)
	return sql
}

func (sd *SQLGoDelete) SetParams(params ...interface{}) *SQLGoDelete {
	if len(params) < 1 {
		return sd
	}
	sd.params = append(sd.params, params...)
	return sd
}

func (sf *SQLGoDelete) GetParams() []interface{} {
	return sf.params
}

func (sf *SQLGoDelete) SetParamsCount(paramsCount int) *SQLGoDelete {
	sf.paramCount = paramsCount
	return sf
}

func (sf *SQLGoDelete) GetParamsCount() int {
	return sf.paramCount
}
