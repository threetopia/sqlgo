package sqlgo

import "fmt"

type SQLGoSelect interface {
	SQLSelect(values ...sqlGoSelectValue) SQLGoSelect
	SetSQLSelect(value sqlGoValue, alias sqlGoAlias) SQLGoSelect
	SetSQLSelectTsRank(column sqlGoColumn, lang string, value sqlGoValue, alias sqlGoAlias) SQLGoSelect
	SetSQLSelectDistinct(column sqlGoColumn) SQLGoSelect
	SetSQLSelectEmbedding(prefix string, column sqlGoColumn, operator sqlGoOperator, value sqlGoValue, alias sqlGoAlias) SQLGoSelect
	SetSQLSelectEmbeddingArray(prefix string, column sqlGoColumn, operator sqlGoOperator, value []sqlGoValue, alias sqlGoAlias) SQLGoSelect

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

	sqlGoSelectDistinct struct {
		column sqlGoColumn
	}

	sqlGoSelectEmbedding struct {
		prefix   string
		column   sqlGoColumn
		operator sqlGoOperator
		value    sqlGoValue
	}

	sqlGoSelectEmbeddingArray struct {
		prefix     string
		column     sqlGoColumn
		operator   sqlGoOperator
		valueArray []sqlGoValue
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

func SetSQLSelectDistinct(column sqlGoColumn) sqlGoSelectValue {
	return sqlGoSelectValue{
		value: sqlGoSelectDistinct{
			column: column,
		},
	}
}

func SetSQLSelectEmbedding(prefix string, column sqlGoColumn, operator sqlGoOperator, value sqlGoValue, alias sqlGoAlias) sqlGoSelectValue {
	return sqlGoSelectValue{
		alias: alias,
		value: sqlGoSelectEmbedding{
			prefix:   prefix,
			column:   column,
			operator: operator,
			value:    value,
		},
	}
}

func SetSQLSelectEmbeddingArray(prefix string, column sqlGoColumn, operator sqlGoOperator, valueArray []sqlGoValue, alias sqlGoAlias) sqlGoSelectValue {
	return sqlGoSelectValue{
		alias: alias,
		value: sqlGoSelectEmbeddingArray{
			prefix:     prefix,
			column:     column,
			operator:   operator,
			valueArray: valueArray,
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

func (s *sqlGoSelect) SetSQLSelectDistinct(column sqlGoColumn) SQLGoSelect {
	s.values = append(s.values, SetSQLSelectDistinct(column))
	return s
}

func (s *sqlGoSelect) SetSQLSelectEmbedding(prefix string, column sqlGoColumn, operator sqlGoOperator, value sqlGoValue, alias sqlGoAlias) SQLGoSelect {
	s.values = append(s.values, SetSQLSelectEmbedding(prefix, column, operator, value, alias))
	return s
}

func (s *sqlGoSelect) SetSQLSelectEmbeddingArray(prefix string, column sqlGoColumn, operator sqlGoOperator, valueArray []sqlGoValue, alias sqlGoAlias) SQLGoSelect {
	s.values = append(s.values, SetSQLSelectEmbeddingArray(prefix, column, operator, valueArray, alias))
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
		case sqlGoSelectDistinct:
			sql = fmt.Sprintf("%sDISTINCT ON (%s) %s", sql, vType.column, vType.column)
		case sqlGoSelectEmbedding:
			s.GetSQLGoParameter().SetSQLParameter(vType.value)
			if vType.prefix == "" {
				sql = fmt.Sprintf("%s(%s %s %s)", sql, vType.column, vType.operator, s.GetSQLGoParameter().GetSQLParameterSign(vType.value))
			} else {
				sql = fmt.Sprintf("%s %s(%s %s %s)", sql, vType.prefix, vType.column, vType.operator, s.GetSQLGoParameter().GetSQLParameterSign(vType.value))
			}
		case sqlGoSelectEmbeddingArray:
			sqlArray := "ARRAY["
			for i, arrayValue := range vType.valueArray {
				s.GetSQLGoParameter().SetSQLParameter(arrayValue)
				if i > 0 {
					sqlArray = fmt.Sprintf("%s,", sqlArray)
				}
				sqlArray = fmt.Sprintf("%s%s", sqlArray, s.GetSQLGoParameter().GetSQLParameterSign(arrayValue))
			}
			sqlArray = fmt.Sprintf("%s]", sqlArray)
			if vType.prefix == "" {
				sql = fmt.Sprintf("%s(%s %s %s)", sql, vType.column, vType.operator, sqlArray)
			} else {
				sql = fmt.Sprintf("%s%s(%s %s %s)", sql, vType.prefix, vType.column, vType.operator, sqlArray)
			}
		default:
			sql = fmt.Sprintf("%s%s", sql, vType)
		}
		if value.alias != "" {
			sql = fmt.Sprintf("%s AS %s", sql, value.alias)
		}
	}
	return sql
}
