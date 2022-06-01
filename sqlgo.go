package sqlgo

import (
	"fmt"
	"strings"
)

type KnownTypes interface {
	string | []string |
		int | []int |
		float32 | []float32 |
		float64 | []float64
}

type SQLBuilder struct {
	insertClause SQLInsertValue
	updateClause []SQLUpdateValue
	selectClause []SQLSelectValue
	fromClause   SQLFromValue
	joinClause   []SQLJoinValue
	whereClause  []SQLWhereValue
	parameters   []interface{}
	paramCount   int
	isJoinScope  bool
}

type SQLInsertValue struct {
	Table   string
	Columns []string
	Values  [][]interface{}
}

type SQLInsertValues struct {
	Column string
	Value  []interface{}
}
type SQLUpdateValue struct {
	Column string
	Value  interface{}
}
type SQLSelectValue struct {
	Alias string
	Value interface{}
}

type SQLFromValue struct {
	Alias string
	Value interface{}
}

type SQLJoinValue struct {
	JoinType  string      // INNER, LEFT, RIGHT, OUTER
	Value     interface{} // table name, *SQLBuilder for sub query
	Alias     string
	JoinWhere []SQLWhereValue
}

type SQLWhereValue struct {
	InJoin      bool
	WhereType   string      // AND, OR
	WhereColumn string      // columnName
	Operator    string      // =, <>, IS
	Value       interface{} // anyValue, *SQLBuilder for sub query
}

func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{}
}

func SetInsert(column string, value ...interface{}) SQLInsertValues {
	return SQLInsertValues{
		Column: column,
		Value:  value,
	}
}

func SetColumns(columns ...string) []string {
	return columns
}

func SetValues(values ...interface{}) []interface{} {
	return values
}

func SetUpdate[V string | *SQLBuilder](column string, value V) SQLUpdateValue {
	return SQLUpdateValue{
		Column: column,
		Value:  value,
	}
}

func SetSelect[V string | *SQLBuilder](value V, alias string) SQLSelectValue {
	return SQLSelectValue{
		Alias: alias,
		Value: value,
	}
}

func SetFrom[V string | *SQLBuilder](value V, alias string) SQLFromValue {
	return SQLFromValue{
		Alias: alias,
		Value: value,
	}
}

func SetJoin[V string | *SQLBuilder](joinType string, value V, alias string, JoinWhere ...SQLWhereValue) SQLJoinValue {
	return SQLJoinValue{
		JoinType:  joinType,
		Value:     value,
		Alias:     alias,
		JoinWhere: JoinWhere,
	}
}

func SetWhere[V KnownTypes | *SQLBuilder](whereType string, column string, operator string, value V) SQLWhereValue {
	return SQLWhereValue{
		WhereType:   whereType,
		WhereColumn: column,
		Operator:    operator,
		Value:       value,
	}
}

func (sb *SQLBuilder) SQLInsert(table string, columns []string, values ...[]interface{}) *SQLBuilder {
	sb.insertClause.Table = table
	sb.insertClause.Columns = columns
	sb.insertClause.Values = append(sb.insertClause.Values, values...)
	return sb
}

func (sb *SQLBuilder) SQLSelect(values ...SQLSelectValue) *SQLBuilder {
	sb.selectClause = append(sb.selectClause, values...)
	return sb
}

func (sb *SQLBuilder) SQLFrom(v SQLFromValue) *SQLBuilder {
	sb.fromClause = v
	return sb
}

func (sb *SQLBuilder) setJoinScope() *SQLBuilder {
	sb.isJoinScope = true
	return sb
}

func (sb *SQLBuilder) SQLJoin(values ...SQLJoinValue) *SQLBuilder {
	sb.joinClause = append(sb.joinClause, values...)
	return sb
}

var specialOperator = map[string]string{
	"ANY": "= ANY ",
	"IN":  " IN ",
}

func (sb *SQLBuilder) SQLWhere(values ...SQLWhereValue) *SQLBuilder {
	sb.whereClause = append(sb.whereClause, values...)
	return sb
}

func (sb *SQLBuilder) BuildSQL() string {
	sql := ""
	if sqlInsert := sb.buildSQLInsert(); sqlInsert != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlInsert)
	}
	if sqlSelect := sb.buildSQLSelect(); sqlSelect != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlSelect)
	}
	if sqlFrom := sb.buildSQLFrom(); sqlFrom != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlFrom)
	}
	if sqlJoin := sb.buildSQLJoin(); sqlJoin != "" {
		sql = fmt.Sprintf("%s%s", sql, sqlJoin)
	}
	if sqlWhere := sb.buildSQLWhere(); sqlWhere != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlWhere)
	}
	return sql
}

func (sb *SQLBuilder) setParamCount(count int) *SQLBuilder {
	sb.paramCount = count
	return sb
}
func (sb *SQLBuilder) getParamCount() int {
	return sb.paramCount
}

func (sb *SQLBuilder) SetParams(params ...interface{}) *SQLBuilder {
	sb.parameters = append(sb.parameters, params...)
	return sb
}

func (sb *SQLBuilder) GetParams() []interface{} {
	return sb.parameters
}

func (sb *SQLBuilder) buildSQLInsert() string {
	if sb.insertClause.Values == nil {
		return ""
	}

	sql := fmt.Sprintf("INSERT %s (", sb.insertClause.Table)
	for iColumn, vColumn := range sb.insertClause.Columns {
		if iColumn > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		sql = fmt.Sprintf("%s%s", sql, vColumn)
	}
	sql = fmt.Sprintf("%s)", sql)

	for iValues, vValues := range sb.insertClause.Values {
		if iValues > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		} else {
			sql = fmt.Sprintf("%s VALUES ", sql)
		}

		sqlValues := "("
		for iValue, vValue := range vValues {
			if iValue > 0 {
				sqlValues = fmt.Sprintf("%s, ", sqlValues)
			}
			sqlValues = fmt.Sprintf("%s%s", sqlValues, vValue)
		}
		sqlValues = fmt.Sprintf("%s)", sqlValues)
		sql = fmt.Sprintf("%s%s", sql, sqlValues)
	}
	return sql
}

func (sb *SQLBuilder) buildSQLSelect() string {
	if sb.selectClause == nil {
		return ""
	}

	sql := "SELECT "
	for iSelect, vSelect := range sb.selectClause {
		if iSelect > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}

		alias := ""
		if vSelect.Alias != "" {
			alias = fmt.Sprintf(" AS %s", vSelect.Alias)
		}

		switch vType := vSelect.Value.(type) {
		case string:
			sql = fmt.Sprintf("%s%s%s", sql, vType, alias)
		case *SQLBuilder:
			sql = fmt.Sprintf("%s(%s)%s", sql, vType.setParamCount(sb.getParamCount()).BuildSQL(), alias)
			sb.SetParams(vType.GetParams()...)
			sb.setParamCount(vType.getParamCount())
		}
	}
	return sql
}

func (sb *SQLBuilder) buildSQLFrom() string {
	if sb.fromClause.Value == nil {
		return ""
	}

	alias := ""
	if sb.fromClause.Alias != "" {
		alias = fmt.Sprintf(" AS %s", sb.fromClause.Alias)
	}

	sql := "FROM "
	switch vType := sb.fromClause.Value.(type) {
	case string:
		sql = fmt.Sprintf("%s%s%s", sql, vType, alias)
	case *SQLBuilder:
		sql = fmt.Sprintf("%s(%s)%s", sql, vType.setParamCount(sb.getParamCount()).BuildSQL(), alias)
		sb.SetParams(vType.GetParams()...)
		sb.setParamCount(vType.getParamCount())
	}
	return sql
}

func (sb *SQLBuilder) buildSQLJoin() string {
	if sb.joinClause == nil {
		return ""
	}

	sql := ""
	for _, vJoin := range sb.joinClause {
		switch vType := vJoin.Value.(type) {
		case *SQLBuilder:
			sqlWhere := NewSQLBuilder().setJoinScope().SQLWhere(vJoin.JoinWhere...)
			sql = fmt.Sprintf("%s %s JOIN (%s) AS %s%s",
				sql,
				strings.ToUpper(vJoin.JoinType),
				vType.setParamCount(sb.getParamCount()).BuildSQL(),
				vJoin.Alias,
				strings.ReplaceAll(sqlWhere.setParamCount(sb.getParamCount()).setParamCount(vType.getParamCount()).BuildSQL(), "WHERE", "ON"),
			)
			sb.SetParams(sqlWhere.parameters...).
				setParamCount(sqlWhere.getParamCount()).
				SetParams(vType.GetParams()...).
				setParamCount(vType.getParamCount())
		case string:
			sqlWhere := NewSQLBuilder().setJoinScope().SQLWhere(vJoin.JoinWhere...)
			sql = fmt.Sprintf("%s %s JOIN %s AS %s%s",
				sql,
				strings.ToUpper(vJoin.JoinType),
				vType,
				vJoin.Alias,
				strings.ReplaceAll(sqlWhere.setParamCount(sb.getParamCount()).BuildSQL(), "WHERE", "ON"),
			)
			sb.SetParams(sqlWhere.parameters...)
			sb.setParamCount(sqlWhere.getParamCount())
		default:
			continue
		}
	}

	return sql
}

func (sb *SQLBuilder) buildSQLWhere() string {
	if sb.whereClause == nil {
		return ""
	}

	sql := "WHERE"
	for iWhere, vWhere := range sb.whereClause {
		whereType := " "
		if iWhere > 0 {
			whereType = fmt.Sprintf(" %s ", strings.ToUpper(vWhere.WhereType))
		}

		operator := vWhere.Operator
		if vo, ok := specialOperator[vWhere.Operator]; ok {
			vWhere.Operator = vo
		}

		switch vType := vWhere.Value.(type) {
		case *SQLBuilder:
			sql = fmt.Sprintf("%s%s%s%s(%s)", sql, whereType, vWhere.WhereColumn, vWhere.Operator, vType.setParamCount(sb.getParamCount()).BuildSQL())
			sb.SetParams(vType.parameters...)
			sb.setParamCount(vType.getParamCount() + sb.getParamCount())
		case []string:
			sql = buildWhereSlice(sb, sql, operator, whereType, vWhere, vType)
		case []int:
			sql = buildWhereSlice(sb, sql, operator, whereType, vWhere, vType)
		case []int64:
			sql = buildWhereSlice(sb, sql, operator, whereType, vWhere, vType)
		case []float64:
			sql = buildWhereSlice(sb, sql, operator, whereType, vWhere, vType)
		default:
			if sb.isJoinScope {
				sql = fmt.Sprintf("%s%s%s%s%s", sql, whereType, vWhere.WhereColumn, vWhere.Operator, vType)
			} else {
				sb.SetParams(vType)
				sb.setParamCount(sb.getParamCount() + 1)
				sql = fmt.Sprintf("%s%s%s%s$%d", sql, whereType, vWhere.WhereColumn, vWhere.Operator, sb.getParamCount())
			}
		}
	}

	return sql
}

func buildWhereSlice[V string | int | int64 | float32 | float64](sb *SQLBuilder, sql string, operator string, whereType string, vWhere SQLWhereValue, vType []V) string {
	if operator == "IN" {
		sql = fmt.Sprintf("%s%s%s%s(", sql, whereType, vWhere.WhereColumn, vWhere.Operator)
		for iIn, vIn := range vType {
			delimiter := ""
			if iIn > 0 {
				delimiter = ","
			}
			if sb.isJoinScope {
				sql = fmt.Sprintf("%s%s%x", sql, delimiter, vIn)
			} else {
				sb.SetParams(vIn)
				sb.setParamCount(sb.getParamCount() + 1)
				sql = fmt.Sprintf("%s%s$%d", sql, delimiter, sb.getParamCount())
			}
		}
		sql = fmt.Sprintf("%s)", sql)
	} else {
		sb.SetParams(vType)
		sb.setParamCount(sb.getParamCount() + 1)
		sql = fmt.Sprintf("%s%s%s%s($%d)", sql, whereType, vWhere.WhereColumn, vWhere.Operator, sb.getParamCount())
	}
	return sql
}
