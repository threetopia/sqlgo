package sqlgo

import (
	"fmt"
	"strings"
)

type SQLBuilder struct {
	selectClause []SQLSelectValue
	fromClause   SQLFromValue
	joinClause   []SQLJoinValue
	whereClause  []SQLWhereValue
	parameters   []interface{}
	paramCount   int
	isJoinScope  bool
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

func SetWhere[V string | []string | int | []int | int64 | []int64 | *SQLBuilder](whereType string, column string, operator string, value V) SQLWhereValue {
	return SQLWhereValue{
		WhereType:   whereType,
		WhereColumn: column,
		Operator:    operator,
		Value:       value,
	}
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
	sql := sb.buildSQLSelect()
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

func (sb *SQLBuilder) buildSQLSelect() string {
	if sb.selectClause == nil {
		return ""
	}

	sql := "SELECT "
	for iSelect, vSelect := range sb.selectClause {
		if iSelect > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}

		switch vType := vSelect.Value.(type) {
		case string:
			sql = fmt.Sprintf("%s%s AS %s", sql, vType, vSelect.Alias)
		case *SQLBuilder:
			sql = fmt.Sprintf("%s(%s) AS %s", sql, vType.setParamCount(sb.getParamCount()).BuildSQL(), vSelect.Alias)
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

	sql := "FROM "
	switch vType := sb.fromClause.Value.(type) {
	case string:
		sql = fmt.Sprintf("%s%s AS %s", sql, vType, sb.fromClause.Alias)
	case *SQLBuilder:
		sql = fmt.Sprintf("%s(%s) AS %s", sql, vType.setParamCount(sb.getParamCount()).BuildSQL(), sb.fromClause.Alias)
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
			// sql = fmt.Sprintf("(%s) AS %s", vt.BuildSQL(), v.Alias)
			// sb.Parameters = append(sb.Parameters, vt.Parameters...)
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
		// case string:
		// 	if sb.isJoinScope {
		// 		sql = fmt.Sprintf("%s%s%s%s%s", sql, whereType, vWhere.WhereColumn, vWhere.Operator, vType)
		// 	} else {
		// 		sb.SetParams(vType)
		// 		sb.setParamCount(sb.getParamCount() + 1)
		// 		sql = fmt.Sprintf("%s%s%s%s$%d", sql, whereType, vWhere.WhereColumn, vWhere.Operator, sb.getParamCount())
		// 	}
		// case int:
		// 	if sb.isJoinScope {
		// 		sql = fmt.Sprintf("%s%s%s%s%s", sql, whereType, vWhere.WhereColumn, vWhere.Operator, vType)
		// 	} else {
		// 		sb.SetParams(vType)
		// 		sb.setParamCount(sb.getParamCount() + 1)
		// 		sql = fmt.Sprintf("%s%s%s%s$%d", sql, whereType, vWhere.WhereColumn, vWhere.Operator, sb.getParamCount())
		// 	}
		case []string:
			if operator == "IN" {
				sql = fmt.Sprintf("%s%s%s%s(", sql, whereType, vWhere.WhereColumn, vWhere.Operator)
				for iIn, vIn := range vType {
					delimiter := ""
					if iIn > 0 {
						delimiter = ","
					}
					if sb.isJoinScope {
						sql = fmt.Sprintf("%s%s%s", sql, delimiter, vType)
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
		case []int:
			if operator == "IN" {
				sql = fmt.Sprintf("%s%s%s%s(", sql, whereType, vWhere.WhereColumn, vWhere.Operator)
				for iIn, vIn := range vType {
					delimiter := ""
					if iIn > 0 {
						delimiter = ","
					}
					if sb.isJoinScope {
						sql = fmt.Sprintf("%s%s%d", sql, delimiter, vType)
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
