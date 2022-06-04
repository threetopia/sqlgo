package sqlgo

type SQLGOInsert struct {
	table      string
	columns    []SQLGoInsertColumn
	values     [][]interface{}
	params     []interface{}
	paramCount int
}

type SQLGoInsertColumn string

type SQLGoInsertValue []interface{}

func NewSQLGOInsert() *SQLGOInsert {
	return &SQLGOInsert{}
}

func SetInsColumns(columns ...string) []string {
	return columns
}

func SetInsValues(values ...interface{}) []interface{} {
	return values
}

func (si *SQLGOInsert) SQLInsert(table string, columns []SQLGoInsertColumn, values ...[]interface{}) *SQLGOInsert {
	si.table = table
	si.columns = append(si.columns, columns...)
	si.values = append(si.values, values...)
	return si
}

func (si *SQLGOInsert) SetSQLInsert(table string, columns []SQLGoInsertColumn, values ...[]interface{}) *SQLGOInsert {
	si.table = table
	si.columns = append(si.columns, columns...)
	si.values = append(si.values, values...)
	return si
}

func (si *SQLGOInsert) BuildSQL() string {
	sql := ""
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
