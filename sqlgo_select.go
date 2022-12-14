package sqlgo

import "fmt"

type SQLGoSelect interface {
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoSelect
	GetSQLGoParameter() SQLGoParameter

	SQLSelect(values ...sqlGoSelectValue) SQLGoSelect
	SetSQLSelect(value interface{}, alias sqlGoAlias) SQLGoSelect

	SQLGoMandatory
}

type (
	sqlGoSelect struct {
		values         []sqlGoSelectValue
		sqlGoParameter SQLGoParameter
	}

	sqlGoSelectValue struct {
		alias sqlGoAlias
		value sqlGoValue
	}
)

func NewSQLGoSelect() SQLGoSelect {
	return &sqlGoSelect{
		sqlGoParameter: NewSQLGoParameter(),
	}
}

func SetSQLSelect(value interface{}, alias sqlGoAlias) sqlGoSelectValue {
	return sqlGoSelectValue{
		alias: alias,
		value: value,
	}
}

func (s *sqlGoSelect) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoSelect {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoSelect) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoSelect) SQLSelect(values ...sqlGoSelectValue) SQLGoSelect {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoSelect) SetSQLSelect(value interface{}, alias sqlGoAlias) SQLGoSelect {
	s.values = append(s.values, SetSQLSelect(value, alias))
	return s
}

func (s *sqlGoSelect) BuildSQL() string {
	var sql string
	if len(s.values) < 1 {
		return sql
	}

	sql = "SELECT "
	for i, value := range s.values {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		switch vType := value.value.(type) {
		case SQLGo:
			vType.SetSQLGoParameter(s.GetSQLGoParameter())
			sql = fmt.Sprintf("%s(%s)", sql, vType.BuildSQL())
			s.SetSQLGoParameter(vType.GetSQLGoParameter())
		default:
			sql = fmt.Sprintf("%s%s", sql, vType)
		}
		if value.alias != "" {
			sql = fmt.Sprintf("%s AS %s", sql, value.alias)
		}
	}
	return sql
}
