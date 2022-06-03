package sqlgo

import "fmt"

type SQLGo struct {
	sqlSelect  *SQLGoSelect
	params     []interface{}
	paramCount int
}

func NewSQLGo() *SQLGo {
	return &SQLGo{
		sqlSelect: NewSQLGoSelect(),
	}
}

func (sg *SQLGo) SQLSelect(values ...sqlGoSelectValues) *SQLGo {
	sg.sqlSelect.SQLSelect(values...)
	return sg
}

func (sg *SQLGo) BuildSQL() string {
	sql := ""
	if sqlSelect := sg.sqlSelect.BuildSQL(); sqlSelect != "" {
		sql = fmt.Sprintf("%s%s", sql, sqlSelect)
	}
	return sql
}

func (sg *SQLGo) GetParams() []interface{} {
	return sg.params
}

func (sg *SQLGo) GetParamsCount() int {
	return sg.paramCount
}
