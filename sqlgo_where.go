package sqlgo

type SQLGoWhere interface {
	SQLGoParameter
}

type sqlGoWhere struct {
	SQLGoParameter
}

func NewSQLGoWhere() SQLGoWhere {
	return sqlGoWhere{}
}
