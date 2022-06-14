package sqlgo

import "fmt"

type SQLGo struct {
	sqlInsert      *SQLGOInsert
	sqlUpdate      *SQLGoUpdate
	sqlDelete      *SQLGoDelete
	sqlSelect      *SQLGoSelect
	sqlFrom        *SQLGoFrom
	sqlJoin        *SQLGoJoin
	sqlWhere       *SQLGoWhere
	sqlGroup       *SQLGoGroup
	sqlOffsetLimit *SQLGoOffsetLimit
	params         []interface{}
	paramCount     int
}

func NewSQLGo() *SQLGo {
	return &SQLGo{
		sqlInsert:      NewSQLGOInsert(),
		sqlUpdate:      NewSQLGoUpdate(),
		sqlDelete:      NewSQLGoDelete(),
		sqlSelect:      NewSQLGoSelect(),
		sqlFrom:        NewSQLGoFrom(),
		sqlJoin:        NewSQLGoJoin(),
		sqlWhere:       NewSQLGoWhere(),
		sqlGroup:       NewSQLGoGroup(),
		sqlOffsetLimit: NewSQLGoOffsetLimit(),
	}
}

//
func (sg *SQLGo) SQLInsert(table string, columns []SQLGoInsertColumn, values ...[]SQLGoInsertValue) *SQLGo {
	sg.sqlInsert.SQLInsert(table, columns, values...)
	return sg
}

func (sg *SQLGo) SetSQLInsert(table string) *SQLGo {
	sg.sqlInsert.setSQLInsertTable(table)
	return sg
}

func (sg *SQLGo) SetSQLInsertColumn(columns ...SQLGoInsertColumn) *SQLGo {
	sg.sqlInsert.setSQLInsertColumn(columns...)
	return sg
}

func (sg *SQLGo) SetSQLInsertValue(values ...SQLGoInsertValue) *SQLGo {
	sg.sqlInsert.setSQLInsertValue(SetInsertValues(values...))
	return sg
}

func (sg *SQLGo) SQLUpdate(table string, values ...SQLGoUpdateValue) *SQLGo {
	sg.sqlUpdate.SQLUpdate(table, values...)
	return sg
}

func (sg *SQLGo) SetSQLUpdate(table string) *SQLGo {
	sg.sqlUpdate.setSQLUpdateTable(table)
	return sg
}

func (sg *SQLGo) SetSQLUpdateValue(column string, value interface{}) *SQLGo {
	sg.sqlUpdate.setSQLUpdateValue(SetUpdate(column, value))
	return sg
}

func (sg *SQLGo) SQLDelete(table string) *SQLGo {
	sg.sqlDelete.SQLDelete(table)
	return sg
}

func (sg *SQLGo) SQLSelect(values ...SqlGoSelectValue) *SQLGo {
	sg.sqlSelect.SQLSelect(values...)
	return sg
}

func (sg *SQLGo) SetSQLSelect(value interface{}, alias string) *SQLGo {
	sg.SQLSelect(SetSelect(value, alias))
	return sg
}

func (sg *SQLGo) SQLFrom(table interface{}, alias string) *SQLGo {
	sg.sqlFrom.SQLFrom(table, alias)
	return sg
}

func (sg *SQLGo) SetSQLFrom(table interface{}, alias string) *SQLGo {
	return sg.SQLFrom(table, alias)
}

func (sg *SQLGo) SQLJoin(values ...SQLGoJoinValue) *SQLGo {
	sg.sqlJoin.SQLJoin(values...)
	return sg
}

func (sg *SQLGo) SetSQLJoin(joinType string, table interface{}, alias string, sqlWhere ...SqlGoWhereValue) *SQLGo {
	sg.SQLJoin(SetJoin(joinType, table, alias, sqlWhere...))
	return sg
}

func (sg *SQLGo) SQLWhere(values ...SqlGoWhereValue) *SQLGo {
	sg.sqlWhere.SQLWhere(values...)
	return sg
}

func (sg *SQLGo) SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) *SQLGo {
	sg.SQLWhere(SetWhere(whereType, whereColumn, operator, value))
	return sg
}

func (sg *SQLGo) SQLGroup(columns ...SQLGoGroupColumn) *SQLGo {
	sg.sqlGroup.SQLGroup(columns...)
	return sg
}

func (sg *SQLGo) SetSQLGroup(columns ...SQLGoGroupColumn) *SQLGo {
	sg.SQLGroup(columns...)
	return sg
}

func (sg *SQLGo) SQLOffsetLimit(offset int, limit int) *SQLGo {
	sg.sqlOffsetLimit.SQLOffsetLimit(offset, limit)
	return sg
}

func (sg *SQLGo) SetSQLOffset(offset int) *SQLGo {
	sg.sqlOffsetLimit.SetSQLOffset(offset)
	return sg
}

func (sg *SQLGo) SetSQLLimit(limit int) *SQLGo {
	sg.sqlOffsetLimit.SetSQLLimit(limit)
	return sg
}

func (sg *SQLGo) BuildSQL() string {
	sql := ""
	if sqlInsert := sg.sqlInsert.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlInsert != "" {
		sql = fmt.Sprintf("%s%s", sql, sqlInsert)
		sg.SetParams(sg.sqlInsert.GetParams()...)
		sg.SetParamsCount(sg.sqlInsert.GetParamsCount())
	}
	if sqlUpdate := sg.sqlUpdate.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlUpdate != "" {
		sql = fmt.Sprintf("%s%s", sql, sqlUpdate)
		sg.SetParams(sg.sqlUpdate.GetParams()...)
		sg.SetParamsCount(sg.sqlUpdate.GetParamsCount())
	}
	if sqlDelete := sg.sqlDelete.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlDelete != "" {
		sql = fmt.Sprintf("%s%s", sql, sqlDelete)
		sg.SetParams(sg.sqlDelete.GetParams()...)
		sg.SetParamsCount(sg.sqlDelete.GetParamsCount())
	}
	if sqlSelect := sg.sqlSelect.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlSelect != "" {
		sql = fmt.Sprintf("%s%s", sql, sqlSelect)
		sg.SetParams(sg.sqlSelect.GetParams()...)
		sg.SetParamsCount(sg.sqlSelect.GetParamsCount())
	}
	if sqlFrom := sg.sqlFrom.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlFrom != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlFrom)
		sg.SetParams(sg.sqlFrom.GetParams()...)
		sg.SetParamsCount(sg.sqlFrom.GetParamsCount())
	}
	if sqlJoin := sg.sqlJoin.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlJoin != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlJoin)
		sg.SetParams(sg.sqlJoin.GetParams()...)
		sg.SetParamsCount(sg.sqlJoin.GetParamsCount())
	}
	if sqlWhere := sg.sqlWhere.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlWhere != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlWhere)
		sg.SetParams(sg.sqlWhere.GetParams()...)
		sg.SetParamsCount(sg.sqlWhere.GetParamsCount())
	}
	if sqlGroup := sg.sqlGroup.BuildSQL(); sqlGroup != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlGroup)
	}
	if sqlOffsetLimit := sg.sqlOffsetLimit.BuildSQL(); sqlOffsetLimit != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlOffsetLimit)
	}
	return sql
}

func (sg *SQLGo) SetParams(params ...interface{}) *SQLGo {
	if len(params) < 1 {
		return sg
	}
	sg.params = append(sg.params, params...)
	return sg
}

func (sg *SQLGo) GetParams() []interface{} {
	return sg.params
}

func (sg *SQLGo) SetParamsCount(paramsCount int) *SQLGo {
	sg.paramCount = paramsCount
	return sg
}

func (sg *SQLGo) GetParamsCount() int {
	return sg.paramCount
}
