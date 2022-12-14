package sqlgo

import (
	"fmt"
	"reflect"
)

type SQLGoSelect interface {
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoSelect
	SQLSelect(values ...sqlGoSelectValue) SQLGoSelect
	SetSQLSelect(alias sqlGoAlias, value interface{}) SQLGoSelect
	SQLGoMandatory
}

type (
	sqlGoSelect struct {
		parameter SQLGoParameter
		values    []sqlGoSelectValue
	}

	sqlGoSelectValue struct {
		alias sqlGoAlias
		value sqlGoValue
	}
)

func NewSQLGoSelect() SQLGoSelect {
	return &sqlGoSelect{
		parameter: NewSQLGoParameter(),
	}
}

func SetSQLSelect(alias sqlGoAlias, value interface{}) sqlGoSelectValue {
	return sqlGoSelectValue{
		alias: alias,
		value: value,
	}
}

func (s *sqlGoSelect) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoSelect {
	s.parameter = sqlGoParameter
	return s
}

func (s *sqlGoSelect) SQLSelect(values ...sqlGoSelectValue) SQLGoSelect {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoSelect) SetSQLSelect(alias sqlGoAlias, value interface{}) SQLGoSelect {
	s.values = append(s.values, SetSQLSelect(alias, value))
	return s
}

func (s *sqlGoSelect) BuildSQL() string {
	sql := "SELECT "
	if len(s.values) > 0 {
		for i, value := range s.values {
			if i > 0 {
				sql = fmt.Sprintf("%s, ", sql)
			}
			xType := reflect.TypeOf(value.value)
			fmt.Println("=================================================", xType)
			switch vType := value.value.(type) {
			case SQLGo:
				// vType.SetSQLGoParameter()
				sql = fmt.Sprintf("%s(%s)", sql, vType.BuildSQL())
			default:
				sql = fmt.Sprintf("%s%s", sql, vType)
			}
			if value.alias != "" {
				sql = fmt.Sprintf("%s AS %s", sql, value.alias)
			}
		}
	}
	return sql
}
