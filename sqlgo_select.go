package sqlgo

import "fmt"

type SQLGoSelect interface {
	SQLSelect(values ...sqlGoSelectValue) SQLGoSelect
	SetSQLSelect(value sqlGoValue, alias sqlGoAlias) SQLGoSelect
	SetSQLSelectTsRank(column sqlGoColumn, lang string, value sqlGoValue, alias sqlGoAlias) SQLGoSelect

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoSelect
	SQLGoBase
}

type (
	sqlGoSelect struct {
		values         sqlGoSelectValues
		sqlGoParameter SQLGoParameter
	}

	sqlGoSelectValue struct {
		alias sqlGoAlias
		value sqlGoValue
	}

	sqlGoSelectValues []sqlGoSelectValue

	sqlGoSelectTsRank struct {
		column sqlGoColumn
		lang   string
		value  sqlGoValue
	}
)

func NewSQLGoSelect() SQLGoSelect {
	return new(sqlGoSelect)
}

func SetSQLSelect(value sqlGoValue, alias sqlGoAlias) sqlGoSelectValue {
	return sqlGoSelectValue{
		alias: alias,
		value: value,
	}
}

func SetSQLSelectTsRank(column sqlGoColumn, lang string, value sqlGoValue, alias sqlGoAlias) sqlGoSelectValue {
	return sqlGoSelectValue{
		alias: alias,
		value: sqlGoSelectTsRank{
			column: column,
			lang:   lang,
			value:  value,
		},
	}
}

func (s *sqlGoSelect) SQLSelect(values ...sqlGoSelectValue) SQLGoSelect {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoSelect) SetSQLSelect(value sqlGoValue, alias sqlGoAlias) SQLGoSelect {
	s.values = append(s.values, SetSQLSelect(value, alias))
	return s
}

func (s *sqlGoSelect) SetSQLSelectTsRank(column sqlGoColumn, lang string, value sqlGoValue, alias sqlGoAlias) SQLGoSelect {
	s.values = append(s.values, SetSQLSelectTsRank(column, lang, value, alias))
	return s
}

func (s *sqlGoSelect) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoSelect {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoSelect) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
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
		case sqlGoSelectTsRank:
			s.GetSQLGoParameter().SetSQLParameter(vType.value)
			sql = fmt.Sprintf("%sts_rank(%s, to_tsquery('%s', %s))", sql, vType.column, vType.lang, s.GetSQLGoParameter().GetSQLParameterSign(vType.value))
		default:
			sql = fmt.Sprintf("%s%s", sql, vType)
		}
		if value.alias != "" {
			sql = fmt.Sprintf("%s AS %s", sql, value.alias)
		}
	}
	return sql
}
