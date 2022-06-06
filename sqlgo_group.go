package sqlgo

import "fmt"

type SQLGoGroup struct {
	columns []SQLGoGroupColumn
}

type SQLGoGroupColumn string

func NewSQLGoGroup() *SQLGoGroup {
	return &SQLGoGroup{}
}

func (sg *SQLGoGroup) SQLGroup(columns ...SQLGoGroupColumn) *SQLGoGroup {
	sg.columns = append(sg.columns, columns...)
	return sg
}

func (sg *SQLGoGroup) BuildSQL() string {
	if len(sg.columns) < 1 {
		return ""
	}

	sql := "GROUP BY "
	for i, v := range sg.columns {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		sql = fmt.Sprintf("%s%s", sql, v)
	}
	return sql
}
