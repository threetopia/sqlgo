package sqlgo

import (
	"fmt"
	"strings"
)

type SQLBuilder struct {
	selectClause []string
	fromClause   string
	joinClause   []string
	whereClause  []string
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
	for _, v := range values {
		sql := ""
		switch vt := v.Value.(type) {
		case *SQLBuilder:
			sb.SetParams(vt.parameters...)
			sb.setParamCount(vt.getParamCount() + sb.getParamCount())
			fmt.Println("========== SQLSelect.getParamCount =========", sb.getParamCount())
			sql = fmt.Sprintf("(%s) AS %s", vt.BuildSQL(), v.Alias)
		case string:
			sql = fmt.Sprintf("%s AS %s", vt, v.Alias)
		default:
			continue
		}
		sb.selectClause = append(sb.selectClause, sql)
	}
	return sb
}

func (sb *SQLBuilder) SQLFrom(v SQLFromValue) *SQLBuilder {
	sql := ""
	switch vt := v.Value.(type) {
	case *SQLBuilder:
		sb.SetParams(vt.parameters...)
		sb.setParamCount(vt.getParamCount() + sb.getParamCount())
		fmt.Println("========== SQLFrom.getParamCount =========", sb.getParamCount())
		sql = fmt.Sprintf("(%s) AS %s", vt.BuildSQL(), v.Alias)
	case string:
		sql = fmt.Sprintf("%s AS %s", vt, v.Alias)
	}
	sb.fromClause = sql

	return sb
}

func (sb *SQLBuilder) setJoinScope() *SQLBuilder {
	sb.isJoinScope = true
	return sb
}

func (sb *SQLBuilder) SQLJoin(values ...SQLJoinValue) *SQLBuilder {
	for _, v := range values {
		sql := ""
		switch vt := v.Value.(type) {
		case *SQLBuilder:
			// sql = fmt.Sprintf("(%s) AS %s", vt.BuildSQL(), v.Alias)
			// sb.Parameters = append(sb.Parameters, vt.Parameters...)
		case string:
			sqlWhere := NewSQLBuilder().setJoinScope().SQLWhere(v.JoinWhere...)
			sb.SetParams(sqlWhere.parameters...)
			sb.setParamCount(sqlWhere.getParamCount() + sb.getParamCount())
			fmt.Println("========== SQLJoin.getParamCount =========", sb.getParamCount())
			sql = fmt.Sprintf(" %s JOIN %s AS %s%s", strings.ToUpper(v.JoinType), vt, v.Alias, sqlWhere.BuildSQL())
		default:
			continue
		}
		sb.joinClause = append(sb.joinClause, sql)
	}
	return sb
}

var specialOperator = map[string]string{
	"ANY": "= ANY ",
	"IN":  " IN ",
}

func (sb *SQLBuilder) SQLWhere(values ...SQLWhereValue) *SQLBuilder {
	for i, v := range values {
		whereType := " "
		if i > 0 {
			whereType = fmt.Sprintf(" %s ", strings.ToUpper(v.WhereType))
		}

		sql := ""
		operator := v.Operator
		if vo, ok := specialOperator[v.Operator]; ok {
			v.Operator = vo
		}

		switch vt := v.Value.(type) {
		case *SQLBuilder:
			sb.SetParams(vt.parameters...)
			sb.setParamCount(vt.getParamCount() + sb.getParamCount())
			sql = fmt.Sprintf("%s%s%s(%s)", whereType, v.WhereColumn, v.Operator, vt.BuildSQL())
			fmt.Println("========== SQLWhere.SQLBuilder.getParamCount =========", sb.getParamCount())
		case string:
			if !sb.isJoinScope {
				sb.SetParams(vt)
				sb.setParamCount(sb.getParamCount() + 1)
				fmt.Println("========== SQLWhere.string.getParamCount =========", sb.getParamCount())
				sql = fmt.Sprintf("%s%s%s$%d", whereType, v.WhereColumn, v.Operator, sb.getParamCount())
			} else {
				sql = fmt.Sprintf("%s%s%s%s", whereType, v.WhereColumn, v.Operator, vt)
			}
		case []string:
			if operator == "IN" {
				sql = fmt.Sprintf("%s%s%s(", whereType, v.WhereColumn, v.Operator)
				for iIn, vIn := range vt {
					delimiter := ""
					if iIn > 0 {
						delimiter = ","
					}
					if !sb.isJoinScope {
						sb.SetParams(vIn)
						sb.setParamCount(sb.getParamCount() + 1)
						fmt.Println("========== SQLWhere.[]string.getParamCount =========", sb.getParamCount())
						sql = fmt.Sprintf("%s%s$%d", sql, delimiter, sb.getParamCount())
					} else {
						sql = fmt.Sprintf("%s%s%s", sql, delimiter, vt)
					}
				}
				sql = fmt.Sprintf("%s)", sql)
			} else {
				sb.SetParams(vt)
				sb.setParamCount(sb.getParamCount() + 1)
				fmt.Println("========== SQLWhere.else.getParamCount =========", sb.getParamCount())
				sql = fmt.Sprintf("%s%s%s($%d)", whereType, v.WhereColumn, v.Operator, sb.getParamCount())
			}
		default:
			continue
		}
		sb.whereClause = append(sb.whereClause, sql)
	}
	return sb
}

func (sb *SQLBuilder) BuildSQL() string {
	sql := ""
	selectDelimiter := ""
	for i, v := range sb.selectClause {
		if i < 1 {
			sql = fmt.Sprintf("SELECT %s", sql)
		} else {
			selectDelimiter = ", "
		}
		sql = fmt.Sprintf("%s%s%s", sql, selectDelimiter, v)
	}

	if sb.fromClause != "" {
		sql = fmt.Sprintf("%s FROM %s", sql, sb.fromClause)
	}

	for _, v := range sb.joinClause {
		v = strings.ReplaceAll(v, " WHERE ", " ON ")
		sql = fmt.Sprintf("%s%s%s", sql, "", v)
	}

	for i, v := range sb.whereClause {
		if i < 1 {
			sql = fmt.Sprintf("%s WHERE", sql)
		}
		sql = fmt.Sprintf("%s%s", sql, v)
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
